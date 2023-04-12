package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
	"net/http"
	"time"
)

var Terminal terminal

type terminal struct{}

func (t *terminal) WsHandler(w http.ResponseWriter, r *http.Request) {
	//解析form入参，其实是get请求，获取相关参数
	if err := r.ParseForm(); err != nil {
		return
	}
	namespace := r.Form.Get("namespace")
	podName := r.Form.Get("pod_name")
	containerName := r.Form.Get("container_name")
	cluster := r.Form.Get("cluster")
	logger.Info("exec pod: %s, container: %s, namespace: %s, cluster: %s \n", podName, containerName, namespace, cluster)
	//获取集群的client
	client, err := K8s.GetClient(cluster)
	if err != nil {
		return
	}
	//加载k8s配置
	conf, err := clientcmd.BuildConfigFromFlags("", K8s.KubeConfMap[cluster])
	if err != nil {
		return
	}
	//实例化TerminalSession
	pty, err := NewTerminalSession(w, r, nil)
	if err != nil {
		logger.Error("get pty failed: %v\n", err)
		return
	}
	//处理关闭
	defer func() {
		logger.Info("close session.")
		pty.Close()
	}()
	//定义POST请求的连接内容
	//https:"192.168.1.11:6443/api/v1/namespaces/default/pods/nginx-wf2-778d88d7c-
	//7rmsk/exec?command=%2Fbin%2Fbash&container=nginxwf2&stderr=true&stdin=true&stdout=true&tty=true
	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: containerName,
			Command:   []string{"/bin/bash"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	logger.Info(req.URL())
	//remotecommand主要实现了http转SPDY协议， 添加X-Stream-Protocol-Version相关的header
	executor, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		return
	}
	//建立连接之后从请求的stream中发送、读取数据
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		TerminalSizeQueue: pty,
		Tty:               true,
	})
	if err != nil {
		msg := fmt.Sprintf("Exec to pod error: %v \n", err)
		logger.Error(msg)
		//将报错发送给web终端，给用户看
		pty.Write([]byte(msg))
		//触发websocket的关闭
		pty.Done()
	}
}

//TerminalSession 定义TerminalSession结构体
//wsConn是websocket连接
//sizeChan用来定义终端的宽和高
//doneChan用于标记退出终端
type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

//TerminalMessage 定义终端交互的内容格式，这个内容格式要遵循remoteCommand规范
//Operation用于定义操作类型
//Data是具体的数据内容
//Rows和Cols也就是终端的行数和列数，也就是宽和高，组成sizeChan
type TerminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

//处理webSocket的协议升级
var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 2
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return upgrader
}()

//NewTerminalSession 负责new一个TerminalSession实例，用于接管输入输出和升级ws协议
func NewTerminalSession(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*TerminalSession, error) {
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	session := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}
	return session, nil
}

//Done 关闭doneChan,关闭后触发终端
func (t *TerminalSession) Done() {
	close(t.doneChan)
}

//Next 定义调整终端的尺寸或退出终端
func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	///监听两个channel，哪个有数据了，去做对应的操作
	//当 t.doneChan 被关闭后，任何试图向该 channel 发送数据的操作都会直接 panic。
	//而对于那些已经从 t.doneChan 中读取或者正在等待从 t.doneChan 中读取数据的操作，它们都将立即返回
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}

//重写输入,输入的对象是web终端，接收web终端输入的内容
func (t *TerminalSession) Read(p []byte) (int, error) {
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		logger.Error(fmt.Sprintf("read message err: %v\n", err))
		return 0, err
	}
	var msg TerminalMessage
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		logger.Error(fmt.Sprintf("read parse message err: %v\n", err))
		return 0, err
	}
	//做一层拦截
	//比如：取消自动补全 日志收集
	//...
	//如果是"stdin"操作，则将消息数据拷贝到给定的字节数组p中并返回拷贝的字节数和nil。
	//如果是"resize"操作，则将消息中的行数和列数发送到t.sizeChan管道中，并返回0和nil。
	//如果是"ping"操作，则直接返回0和nil。
	//否则，记录日志并返回错误。
	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{
			Width:  msg.Cols,
			Height: msg.Rows,
		}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		logger.Error(fmt.Sprintf("unknown message type %s\n", msg.Operation))
		return 0, errors.New(fmt.Sprintf("unknown message type %s\n", msg.Operation))
	}
}

//Write 重写输出，接收web端的指令后，将结果返回出去
func (t *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(TerminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		logger.Error(fmt.Sprintf("write parse message err: %v\n", err))
		return 0, err
	}
	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		logger.Error(fmt.Sprintf("write message err: %v\n", err))
		return 0, err
	}
	return len(p), nil
}

//Close 用于关闭websocket连接
func (t *TerminalSession) Close() error {
	return t.wsConn.Close()
}
