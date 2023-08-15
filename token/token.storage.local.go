package token

// localTokenStorage 本地验证token
type localTokenStorage struct {
	Provider *JwtProvider
}

func NewLocalStorage(provider *JwtProvider) (s Storage) {
	s = &localTokenStorage{
		Provider: provider,
	}
	return
}

func (lts *localTokenStorage) Verify(key string) (p JwtPayload, err error) {
	claim, err := lts.Provider.Verify([]byte(key))
	if err != nil {
		return
	}
	err = claim.Claims(&p)
	return
}
