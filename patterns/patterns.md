# ゴルーチンやチャネルを用いた実装パターン

1. 拘束
```
func restFunc() <-chan int {
	// 1. チャネルを定義
	result := make(chan int)

	// 2. ゴールーチンを立てて
	go func() {
		defer close(result) // 4. closeするのを忘れずに

		// 3. その中で、resultチャネルに値を送る処理をする
		// (例)
		for i := 0; i < 5; i++ {
			result <- 1
		}
	}()

	// 5. 返り値にresultチャネルを返す
	return result
}
```
 resultチャネルが使えるスコープをrestFunc内に留める(=拘束する)ことで、
あらぬところから送信が行われないように保護することができ、安全性が高まります。
ちなみに、resultのように、特定の値をひたすら生成するだけのチャネルをジェネレータという。

2. select分
```
select {
case num := <-gen1:  // gen1チャネルから受信できるとき
	fmt.Println(num)
case num := <-gen2:  // gen2チャネルから受信できるとき
	fmt.Println(num)
default:  // どっちも受信できないとき
	fmt.Println("neither chan cannot use")
}
```
 受信可能なチャネルから値を得ることができる。
これと同じことをifで実装しようとすると、最初の分岐条件がFalseだった時点でそのゴルーチンがロックされてしまい、デッドロックになる可能性がある。
そのため、以上のようにselect文を用いるようにする。

※ バッファありチャネル
make(chan int, 2): int型のチャネルをサイズ2のバッファ付きで作成
バッファが詰まったときはチャネルへの送信をブロックする。
