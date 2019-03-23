package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"

	"github.com/stretchr/gomniauth"
)

/*
各プロバイダ・モデルを生成する。認証情報は info.go に格納し、使用する。
*/
func setAuthInfo() {
	gomniauth.SetSecurityKey(securityKey)
	gomniauth.WithProviders(
		google.New(
			googleClientId,
			googleClientSecurityKey,
			"http://localhost:8080/auth/callback/google",
		),
		github.New(
			githubClientId,
			githubClientSecurityKey,
			"http://localhost:8080/auth/callback/github",
		),
	)
}

/*
指定されたプロバイダに対して認証を行う。
*/
func authHandler(w http.ResponseWriter, r *http.Request) {

	// urlの分解。区切り文字は "/"
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]        //処理内容。login あるいは callback.
	provider_name := segs[3] //プロバイダの名称。github あるいは google

	switch action {

	/* loginページから遷移した場合 */
	case "login":
		// gomniauthを使用し、プロバイダ・モデルを生成する
		provider, err := gomniauth.Provider(provider_name)
		if err != nil {
			log.Fatalln("プロバイダの取得に失敗しました。")
		}
		// プロバイダ毎の認証ページへのurlを取得する
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("認証ページの取得に失敗しました。")
		}

		// 提供された認証用のページへリダイレクトする
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)

	/* プロバイダの元での認証を終え、callbackされた場合 */
	case "callback":
		// gomniauthを使用し、プロバイダ・モデルを生成する
		provider, err := gomniauth.Provider(provider_name)
		if err != nil {
			log.Fatalln("プロバイダの取得に失敗しました。")
		}
		// 提供されたURLから、認証に必要な情報を抜き出す。
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("認証情報の取得に失敗しました", err)
		}
		// 認証情報を使用して、プロバイダからUserオブジェクトを取得する。
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("ユーザ情報の取得に失敗しました。")
		}

		// プロバイダから提供されたuserオブジェクトより情報を抜き出す。
		authCookieValue := objx.New(map[string]interface{}{
			"name":       user.Name(),      //ユーザ名
			"avatar_url": user.AvatarURL(), //プロフ画像のURL
			"provider":   provider_name,
		}).MustBase64()

		// 抜き出した情報をクッキーに詰める。Nameは"auth"
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/after",
		})

		// ログイン後の画面に遷移する
		w.Header()["Location"] = []string{"/after"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
