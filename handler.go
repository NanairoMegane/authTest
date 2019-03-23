package main

import (
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/objx"
)

/*
ルートディレクトリのアクセスに対し、ユーザのステータスに応じて遷移先を変更する。
*/
func moveHandler(w http.ResponseWriter, r *http.Request) {

	// 認証情報の有無を確認する
	authCookie, _ := r.Cookie("auth")

	if authCookie != nil {
		// 認証が済んでいれば、afterページにアクセスさせる
		w.Header()["Location"] = []string{"/after"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	} else {
		// 認証が済んでいなければ、loginページにアクセスさせる
		w.Header()["Location"] = []string{"/login"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

/*
HTMLテンプレートをサーブするためのハンドラ
*/
type templateHandler struct {
	once     sync.Once          //HTMLテンプレートを１度だけコンパイルするための指定
	filename string             //テンプレートとしてHTMLファイル名を指定
	tmpl     *template.Template //テンプレート
}

/*
templateHandlerをhttp.Handleに適合させるため、ServeHttpを実装する
*/
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// テンプレートディレクトリを指定する
	path, err := filepath.Abs("./templates/")
	if err != nil {
		panic(err)
	}

	// 指定された名称のテンプレートファイルを一度だけコンパイルする
	t.once.Do(
		func() {
			t.tmpl = template.Must(template.ParseFiles(path + t.filename))
		})

	// HTMLに渡すデータ
	data := map[string]interface{}{
		"Host": r.Host,
	}

	// 認証が済んでいれば、クッキーの値を渡す。
	authCookie, err := r.Cookie("auth")
	if err == nil && authCookie != nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	// テンプレートをパースする際、テンプレートに渡すデータも指定する。
	t.tmpl.Execute(w, data)
}

/*
/logout に対するハンドラ。クッキー情報を削除し、login画面へ遷移させる。
*/
func logoutHandler(w http.ResponseWriter, r *http.Request) {

	// 既に設定されている"auth"のクッキー情報を削除する。
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "",
		MaxAge: -1,
	})

	w.Header()["Location"] = []string{"/login"}
	w.WriteHeader(http.StatusTemporaryRedirect)
}
