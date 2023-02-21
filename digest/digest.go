package digest

// SaltingFunc 混淆器
//
// plaintext 明文
type SaltingFunc func(plaintext string) string

// Hash checksum计算器
//
// 自行实现或在 Digests 列表中选择
type Hash func(plain []byte) (digest []byte, err error)

var defaultDigestFunc = sm3Digest

var Digests = map[string]Hash{
	"SM3":    sm3Digest,
	"BCRYPT": bcryptDigest,
	"SHA512": sha512Digest,
}

// Generate 创建密文
//
// plaintext 明文
//
// nanoTime 时间戳 time.UnixNano()
//
// salting SaltingFunc 混淆文生成器，最多只会处理一个，如不提供将使用默认算法
func Generate(plaintext string, nanoTime int64, salting ...SaltingFunc) (r string) {
	r = GenerateWithAlg(plaintext, nanoTime, defaultDigestFunc, salting...)
	return
}

func GenerateWithAlg(plaintext string, nanoTime int64, alg Hash, salting ...SaltingFunc) (r string) {
	r = genPass(plaintext, nanoTime, alg, salting...)
	return
}

func Verify(plaintext string, encodedPassword string, nanoTime int64, salting ...SaltingFunc) (r bool) {
	r = VerifyWithAlg(plaintext, encodedPassword, nanoTime, defaultDigestFunc, salting...)
	return
}

func VerifyWithAlg(plaintext string, encodedPassword string, nanoTime int64, alg Hash, salting ...SaltingFunc) (r bool) {
	encoded := genPass(plaintext, nanoTime, alg, salting...)
	r = encoded == encodedPassword
	return
}
