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

3. セマフォ
```
var sem = make(chan int, MaxOutstanding)

func handle(r *Request) {
    sem <- 1    // Wait for active queue to drain.
    process(r)  // May take a long time.
    <-sem       // Done; enable next request to run.
}

func Serve(queue chan *Request) {
    for {
        req := <-queue
        go handle(req)  // Don't wait for handle to finish.
    }
}
```
 Serveは、queueチャネルからリクエストを受け取って、それをhandleする。
ですが、このままだと際限なくhandle関数を実行するゴールーチンが立ち上がってしまいます。それをセマフォとして制御するのがバッファありのsemチャネルです。
handle関数の中で、リクエストを受け取ったらsemに値を1つ送信、リクエストを処理し終えたらsemから値を1つ受信という操作をしています。
もしもsemチャネルがいっぱいになったら、sem <- 1の実行がブロックされます。そのため、semチャネルの最大バッファ数以上のゴールーチンが立ち上がることを防いでいます。

また、これをうまく利用すれば、リーキーバケットアルゴリズムの実装が容易にできる。

※ バッファありチャネル
make(chan int, 2): int型のチャネルをサイズ2のバッファ付きで作成
バッファが詰まったときはチャネルへの送信をブロックする。


4. メインルーチンからサブルーチンを停止させる
 使ってないゴールーチンを稼働したまま放っておくということは、そのスタック領域をGCされないまま放っておくことになる。
これはパフォーマンス的にあまりよくない事態を引き起こし、このような現象のことをゴールーチンリークという。

```
func generator(done chan struct{}) <-chan int {
	result := make(chan int)
	go func() {
		defer close(result)
		for {
			select {
			case <-done: 
                // doneチャネルは空なので、こっちの処理はブロックされ実行されない。
                // しかし、doneチャネルがcloseされると、ゼロ値のnilが返ってくるので実行される。
				break
			case result <- 1:
			}
		}
	}()
	return result
}

func main() {
	done := make(chan struct{})

	result := generator(done)
	for i := 0; i < 5; i++ {
		fmt.Println(<-result)
	}
	close(done)
}
```
generatorはひたすら無限に1を生成するジェネレータ。doneチャネルを使って終了を制御している。
※ doneチャネルはclose操作のみ行い、値などは特に渡さないため、chan struct{}型で定義しておく。こうするとメモリ領域の節約ができる。

5. FanIn(ファン-イン)
 複数個あるチャネルから受信した値を、1つの受信用チャネルの中にまとめる方法をFanInといいます。
```
func fanIn1(done chan struct{}, c1, c2 <-chan int) <-chan int {
	result := make(chan int)

	go func() {
		defer fmt.Println("closed fanin")
		defer close(result)
		for {
			// caseはfor文で回せないので(=可変長は無理)
			// 統合元のチャネルがスライスでくるとかだとこれはできない
			select {
			case <-done:
				fmt.Println("done")
				return
			case num := <-c1:
				fmt.Println("send 1")
				result <- num
			case num := <-c2:
				fmt.Println("send 2")
				result <- num
			default:
				fmt.Println("continue")
				continue
			}
		}
	}()

	return result
}

func main() {
	done := make(chan struct{})

	gen1 := generator(done, 1) // int 1をひたすら送信するチャネル(doneで止める)
	gen2 := generator(done, 2) // int 2をひたすら送信するチャネル(doneで止める)

	result := fanIn1(done, gen1, gen2) // 1か2を受け取り続けるチャネル
	for i := 0; i < 5; i++ {
		<-result
	}
	close(done)
	fmt.Println("main close done")

	// resultチャネルに値が残ってしまうことがあるので、ゴルーチンリークにならないよう後片付け。
	for {
		if _, ok := <-result; !ok {
			break
		}
	}
}
```

6. FanIn応用
 FanInでまとめたいチャネル群が可変長変数やスライスで与えられている場合は、select文を直接使用することができないため、以下のようにする。
```
func fanIn2(done chan struct{}, cs ...<-chan int) <-chan int {
	result := make(chan int)

	var wg sync.WaitGroup
	wg.Add(len(cs))

	for i, c := range cs {
		// FanInの対象になるチャネルごとに個別にゴールーチンを立てちゃう
		go func(c <-chan int, i int) {
			defer wg.Done()

			for num := range c {
				select {
				case <-done:
					fmt.Println("wg.Done", i)
					return
				case result <- num:
					fmt.Println("send", i)
				}
			}
		}(c, i)
	}

	go func() {
		// selectでdoneが閉じられるのを待つと、個別に立てた全てのゴールーチンを終了できる保証がない
		wg.Wait()
		fmt.Println("closing fanin")
		close(result)
	}()

	return result
}
```

7. タイムアウトの実装
1. time.After関数を使用した場合
 time.After関数は、引数d時間経ったら値を送信するチャネルを返す関数。

- 例「1秒以内にselectできるならずっとそうする、できなかったらタイムアウト」
```
for {
    select {
    case s := <-ch1:
        fmt.Println(s)
    case <-time.After(1 * time.Second): // ch1が受信できないまま1秒で発動
        fmt.Println("time out")
        return
    /*
    }
}
```

- 例「select文を実行し続けるのを1秒間行う」
```
timeout := time.After(1 * time.Second)

// このforループを1秒間ずっと実行し続ける
for {
	select {
	case s := <-ch1:
		fmt.Println(s)
	case <-timeout:
		fmt.Println("time out")
		return
	default:
		fmt.Println("default")
		time.Sleep(time.Millisecond * 100)
	}
}
```

2. time.NewTimerの構造体を使用した場合
time.NewTimerのコード
```
type Timer struct {
	C <-chan Time
	// contains filtered or unexported fields
}

func NewTimer(d Duration) *Timer
```

- 例「1秒以内にselectできるならずっとそうする、できなかったらタイムアウト」
```
for {
	t := time.NewTimer(1 * time.Second)
	defer t.Stop()

	select {
	case s := <-ch1:
		fmt.Println(s)
	case <-t.C:
		fmt.Println("time out")
		return
	}
}
```

- 例「select文を実行し続けるのを1秒間行う」
```
t := time.NewTimer(1 * time.Second)
defer t.Stop()

for {
	select {
	case s := <-ch1:
		fmt.Println(s)
	case <-t.C:
		fmt.Println("time out")
		return
	default:
		fmt.Println("default")
		time.Sleep(time.Millisecond * 100)
	}
}
```

- time.Afterとtime.NewTimerの使い分け
 time.After(d)で得られるものはNewTimer(d).Cと同じで、
内包されているタイマーは、作動されるまでガベージコレクトによって回収されることはない。
効率を重視する場合、time.NewTimerの方を使い、タイマーが不要になったタイミングでStopメソッドを呼んでください。

8. 結果のどれかを使う(moving on)
例: DBへのコネクションConnが複数個存在して、その中から得られた結果のうち一番早く返ってきたものを使って処理をしたいという場合。
```
func Query(conns []Conn, query string) Result {
    ch := make(chan Result, len(conns))
	// connから結果を得る作業を並行実行
    for _, conn := range conns {
        go func(c Conn) {
            select {
            case ch <- c.DoQuery(query):
            default:
            }
        }(conn)
    }
    return <-ch
}

func main() {
	// 一番早くchに送信されたやつだけがここで受け取ることができる
	result := Query(conns, query)
	fmt.Println(result)
}
```
「doneチャネルを使ってのルーチン閉じ作業」は省略しているので注意。
