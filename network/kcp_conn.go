package network

import (
	"github.com/po2656233/goleaf/log"
	"github.com/xtaci/kcp-go"
	"net"
	"sync"
	"time"
)

type KCPConnSet map[*kcp.UDPSession]struct{}

type KCPConn struct {
	sync.Mutex
	conn      *kcp.UDPSession
	writeChan chan []byte
	closeFlag bool
	msgParser *MsgParser
}

func newKCPConn(conn *kcp.UDPSession, pendingWriteNum int, msgParser *MsgParser) *KCPConn {
	kcpConn := new(KCPConn)
	kcpConn.conn = conn
	kcpConn.writeChan = make(chan []byte, pendingWriteNum)
	kcpConn.msgParser = msgParser
	kcpConn.conn.SetWriteDelay(false) // 不延迟
	//nodelay ：是否启用 nodelay模式，0不启用；1启用。
	//interval ：协议内部工作的 interval，单位毫秒，比如 10ms或者 20ms
	//resend ：快速重传模式，默认0关闭，可以设置2（2次ACK跨越将会直接重传）
	//kcpConn.conn.SetNoDelay(0, 40, 0, 0) //普通模式
	kcpConn.conn.SetNoDelay(1, 10, 2, 1) //极速模式

	go func() {
		for b := range kcpConn.writeChan {
			if b == nil {
				break
			}

			_, err := conn.Write(b)
			if err != nil {
				break
			}
		}

		conn.Close()
		kcpConn.Lock()
		kcpConn.closeFlag = true
		kcpConn.Unlock()
	}()

	return kcpConn
}

func (kcpConn *KCPConn) doDestroy() {
	kcpConn.conn.SetDeadline(time.Time{})
	kcpConn.conn.Close()

	if !kcpConn.closeFlag {
		close(kcpConn.writeChan)
		kcpConn.closeFlag = true
	}
}

func (kcpConn *KCPConn) Destroy() {
	kcpConn.Lock()
	defer kcpConn.Unlock()

	kcpConn.doDestroy()
}

func (kcpConn *KCPConn) Close() {
	kcpConn.Lock()
	defer kcpConn.Unlock()
	if kcpConn.closeFlag {
		return
	}

	kcpConn.doWrite(nil)
	kcpConn.closeFlag = true
}

func (kcpConn *KCPConn) doWrite(b []byte) {
	if len(kcpConn.writeChan) == cap(kcpConn.writeChan) {
		log.Debug("close conn: channel full")
		kcpConn.doDestroy()
		return
	}

	kcpConn.writeChan <- b
}

// b must not be modified by the others goroutines
func (kcpConn *KCPConn) Write(b []byte) {
	kcpConn.Lock()
	defer kcpConn.Unlock()
	if kcpConn.closeFlag || b == nil {
		return
	}

	kcpConn.doWrite(b)
}

func (kcpConn *KCPConn) Read(b []byte) (int, error) {
	return kcpConn.conn.Read(b)
}

func (kcpConn *KCPConn) LocalAddr() net.Addr {
	return kcpConn.conn.LocalAddr()
}

func (kcpConn *KCPConn) RemoteAddr() net.Addr {
	return kcpConn.conn.RemoteAddr()
}

func (kcpConn *KCPConn) ReadMsg() ([]byte, error) {

	kcpConn.conn.SetReadDeadline(time.Now().Add(35 * time.Second))
	b, err := kcpConn.msgParser.ReadKcp(kcpConn)
	kcpConn.conn.SetReadDeadline(time.Time{})
	return b, err
}

func (kcpConn *KCPConn) WriteMsg(args ...[]byte) error {
	return kcpConn.msgParser.WriteKcp(kcpConn, args...)
}
