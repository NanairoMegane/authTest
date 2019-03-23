package main

import (
	"log"
	"net/http"
)

func main() {

	/* ルートディレクトリへのアクセスに対し、認証の有無により遷移先をコントロールする */
	/*
		未認証 → /login
		認証済 → /after
	*/
	http.HandleFunc("/", moveHandler)

	/* /loginへのハンドラ(認証プロバイダの選択) */
	http.Handle("/login", &templateHandler{filename: "/login.html"})

	/* 認証用のgomniauthの初期設定 */
	setAuthInfo()

	/* /authへのハンドラ(指定プロバイダでの認証) */
	http.HandleFunc("/auth/", authHandler)

	/* /afterへのハンドラ(認証後のページ・サーブ) */
	http.Handle("/after", &templateHandler{filename: "/after.html"})

	/* /logoutへのハンドラ(認証情報の削除) */
	http.HandleFunc("/logout", logoutHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("Webサーバの起動に失敗しました。")
	}
}
