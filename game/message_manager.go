package game

import (
	"fmt"
	"math/rand"
	"time"
)

type messageType int

const (
	WelcomeMsg = iota
	START_MSG_NEW
	START_MSG_NEW_SIMPLE
	START_MSG_REPEAT
	InstructionMsg
	GOAL_MSG
	GiveUpMsg
	MoveMsg
	LocationMsg
	GameoverMsg
	InvalidActionMsg
	InquirelyMsg

	ItaiMsg
	ButsukaruMsg

	GoodbyMsg

	RepromptMsg1
	RepromptMsg2
	RepromptMsg3
	RepromptMsg4

	CurrentWeather
	CurrentWeather2
	CurrentWeather2NC
	NoCity
	WeatherNotFound
	Tomete
	Arigato
	Sugoine
	NorthPole
)

var (
	messageMap  map[messageType]string
	messageMap2 map[messageType][]string
	rnd         *rand.Rand
)

func init() {

	seed := time.Now().UnixNano()
	rnd = rand.New(rand.NewSource(seed))

	messageMap = make(map[messageType]string)
	messageMap2 = make(map[messageType][]string)

	messageMap[WelcomeMsg] = "世界の天気にようこそ{[500]}！どこの天気が知りたいですか？"

	messageMap[START_MSG_NEW] = "新しい迷路です。大きさは%dかける%d、スタートは左下%s、ゴールは右上%sです。"
	messageMap[START_MSG_NEW_SIMPLE] = "新しい迷路です。スタートは左下%s、ゴールは右上%sです。"
	messageMap[START_MSG_REPEAT] = "スタートに戻ってきました。スタートは%s、ゴールは%sです。"

	messageMap2[InstructionMsg] = []string{"行きたい方向を言って下さい。", "移動する方向を言って下さい。", "行きたい方向を言うと移動できます。"}

	messageMap[GOAL_MSG] = "おめでとう{[500]}！ゴールです{[500]}！移動回数は%d回。確認回数は%d回です。新しい迷路に挑戦するには{[300]}、新しい迷路{[500]}、と言って下さい。"
	messageMap[GiveUpMsg] = "まあそうゆわんともう一回どうですか？"

	messageMap2[MoveMsg] = []string{"%sに進みました。", "%sに移動しました。"}
	messageMap[LocationMsg] = "現在の位置は{[200]}、%s、です。"

	messageMap[GameoverMsg] = "残念{[500]}、ゲームオーバーです。"

	messageMap2[InvalidActionMsg] = []string{"それは無理。", "残念ながらそれはできません。", "それは無理です。", "それはできません。", "残念ながらできません。"}

	messageMap[InquirelyMsg] = "もう一度言ってください。"

	messageMap2[ItaiMsg] = []string{"いたっ{[500]}！", "痛いっ{[500]}！", "あーあ。"}
	messageMap2[ButsukaruMsg] = []string{"壁にぶつかりました。", "そちらには壁があります。", "壁なので進めません。"}

	messageMap[GoodbyMsg] = "さようなら"

	messageMap[RepromptMsg1] = "もう一度挑戦する場合は、やり直す{[500]}、新しい迷路に挑戦する場合は、新しい迷路{[500]}、と言って下さい。"
	messageMap[RepromptMsg2] = "どちらへ行きますか？"
	messageMap2[RepromptMsg2] = []string{"どちらへ行きますか？", "次はどちらへ行きますか？", "次はどうしますか？", "どちらへ行きましょう？", "さあ、次は？", "では、次はどちらへ？"}
	messageMap[RepromptMsg3] = "新しい迷路で遊ぶのなら、新しい迷路{[500]}、もう一度同じ迷路で遊ぶ場合は、同じ迷路{[400]}、と言って下さい。"
	messageMap[RepromptMsg4] = "どうしますか？"

	messageMap[CurrentWeather] = "現在の%sの天気は%s、気温は%d度です。"
	messageMap[CurrentWeather2] = "現在の%s、%sの天気は%s、気温は%d度です。"
	messageMap[CurrentWeather2NC] = "現在の%sの天気は%s、気温は%d度です。"

	messageMap2[WeatherNotFound] = []string{"ごめんなさい、%sの天気はわかりません。", "すいません、%sの天気は知らないんです。", "え、%sですか。申し訳ありませんがそれは無理です。", "%sはちょっと。ごめんなさい。"}

	messageMap2[NoCity] = []string{"ごめんなさい、よく聞き取れませんでした。もう一度言っていただけますか？", "残念ながらお答えできません。", "申し訳ありません。わかりませんでした。"}

	messageMap2[Tomete] = []string{"また呼んでくださいね。", "またのご利用をお待ちしております。", "では、失礼します。", "はい、さようなら。", "はい、ありがとうございました。"}

	messageMap2[Arigato] = []string{"どういたしまして。", "お役に立てて何よりです。", "なんのこれしき。"}

	messageMap2[Sugoine] = []string{"もったいないお言葉ありがとうございます。", "えっ、そんな{[1000]}てへ。", "ありがとう。"}

	messageMap2[NorthPole] = []string{"現在の%sの天気は不明、気温はめっちゃ低いでしょう。", "えっ、%sに行くんですか？", "%sですか？寒いのは間違いないと思いますが詳しいことはわかりません。"}
}

// Assume messageMap[t] always exists
func GetMessage(t messageType, a ...interface{}) string {
	s0, ok := messageMap[t]
	if !ok {
		fmt.Printf("ERROR! message for %d does not exist.", t)
	}
	return fmt.Sprintf(s0, a...)
}

func GetMessage2(t messageType, a ...interface{}) string {
	ss, ok := messageMap2[t]
	if !ok {
		fmt.Printf("ERROR! message for %d does not exist.", t)
	}
	i := rnd.Intn(len(ss))
	s := ss[i]
	return fmt.Sprintf(s, a...)
}

func GetMessage2Random(t messageType, p float64, a ...interface{}) string {
	r := rnd.Float64()
	if r >= p {
		return ""
	}
	ss, ok := messageMap2[t]
	if !ok {
		fmt.Printf("ERROR! message for %d does not exist.", t)
	}
	i := rnd.Intn(len(ss))
	s := ss[i]
	return fmt.Sprintf(s, a...)
}
