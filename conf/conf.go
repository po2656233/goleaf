package conf

var (
	LenStackBuf = 4096

	// log
	LogLevel string
	LogPath  string
	LogFlag  int

	// console
	ConsolePort   int
	ConsolePrompt string = "Leaf# "
	ProfilePath   string

	// cluster
	ListenAddr      string
	ConnAddrs       []string
	PendingWriteNum int

	// kcp
	DSCP  int = 46        // 差分服务代码点
	SALT      = "goleaf"  // 加盐
	Key       = "kcp-key" // 密钥
	Crypt     = "salsa20" // 加密方式 如果此项为空,则SALT|key无效
)
