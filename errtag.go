package errtag

import (
	"runtime"
	"strings"
)

// ErrorTag はコールした関数を階層的に表現したエラータグを返す
// レシーバ関数と通常関数のみ対応している
// TODO: 無名関数対応
//
// e.g. "github.com/foo/bar/infra/db/mysql/user/repository.go" 内のRepository.Get関数からコールした場合
//   output: "infra/db/mysql/user/Repository.Get"
func ErrorTag() string {
	// 完全なFunction名を分解する ["github.com", "foo", "bar", "...", "abc.Func"]
	splitted := strings.Split(fullFunctionName(), "/")

	result := make([]string, 0, len(splitted))
	result = append(result, extractNonRepositoryPathAndFunctionName(splitted)...)
	result = append(result, decorateFunctionName(splitted)...)

	// パス形式に連結して返す
	return strings.Join(result, "/")
}

// fullFunctionName は "github.com/foo/bar/.../abc.Func" のような、完全な関数名を返します
//
// e.g. github.com/foo/bar/infra/db/mysql/user/repository.go内のRepository.Get関数からコールした場合
//   output: "github.com/foo/bar/infra/db/mysql/user.(*Repository).Get"
func fullFunctionName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}

// extractNonRepositoryPathAndFunctionName はリポジトリパス("github.com/foo/bar")とFunction名("abc.Func")以外を抽出します
//
// e.g.
//   input: ["github.com", "foo", "bar", "baz", "abc.(*Def).Ghi"]
//   output: ["baz"]
func extractNonRepositoryPathAndFunctionName(splitted []string) []string {
	// ["github.com", "foo", "bar"] = 3
	const repoPathCount = 3

	l := len(splitted)
	switch {
	// リポジトリパスの必須要件(github.com/hoge/fuga)を満たしていない場合は空で返す
	case l <= 3:
		return []string{}
	default:
		return splitted[repoPathCount : l-1]
	}
}

// decorateFunctionName はFunction名を加工します
//
// e.g.
//   input: ["github.com", "foo", "bar", "baz", "abc.(*Def).Ghi"]
//   output: ["abc", "Def.Ghi"]
//
// レシーバ関数と通常関数で挙動が変わります
// レシーバ関数:
//   src: ["abc", "(*Def)", "Ghi"]
//   dst: ["abc", "Def.Ghi"]
// 通常関数:
//   src: ["abc", "Def"]
//   dst: ["abc.Def"]
func decorateFunctionName(splitted []string) []string {
	const (
		receiverMethodLength = 3
		normalMethodLength   = 2
	)

	fn := splitted[len(splitted)-1] // foo.Func
	// Dot Splitすると長さが変わるのでそこで判別できる
	dotSplitted := strings.Split(fn, ".")

	switch len(dotSplitted) {
	// for Receiver Function
	// ["foo", "(*Bar)", "Baz"]
	case receiverMethodLength:
		trimRcvName := strings.Trim(dotSplitted[1], "()*") // (*Bar) -> Bar
		fName := trimRcvName + "." + dotSplitted[2]        // "Bar.Baz"
		// ["foo", "Bar.Baz"] を返す
		return []string{dotSplitted[0], fName}
	// for Normal Function
	// ["foo", "Bar"]
	default:
		// Dot連結で返す
		return []string{dotSplitted[0] + "." + dotSplitted[1]}
	}
}
