package main

import (
	"bufio"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	knight    *ebiten.Image // 騎士画像
	knightWin *ebiten.Image // 騎士画像（勝利）
	bullet    *ebiten.Image // エフェクト
	ghost     *ebiten.Image // ゴースト画像
	fontFace  font.Face     // フォント
	judgement bool          // 勝ち負け判定
	highScore int           // ハイスコア
	random    float64       // 乱数
)

var mode int = 0            // モード
var score int = 0           // スコア
var otetuki int = 0         // お手付き数
var knightBool = true       // 騎士の描写に関する真理値
var ghostTime float64 = 3   // ゴーストの時間
var elapsedTime float64 = 0 // 経過時間
var now = time.Now()        // now初期化
var sub float64 = 0         // 差分
var knightx float64         // 騎士のｘ座標
var ghostx float64 = 280    // ゴーストのｘ座標

type Game struct{}

// 初期化
func init() {
	rand.Seed(time.Now().Unix()) // 乱数の種
	var err error

	// 画像の定義
	knight, _, err = ebitenutil.NewImageFromFile("image/Knight1.png")
	if err != nil {
		log.Fatal(err)
	}

	knightWin, _, err = ebitenutil.NewImageFromFile("image/Knight2.png")
	if err != nil {
		log.Fatal(err)
	}

	ghost, _, err = ebitenutil.NewImageFromFile("image/ghost1.png")
	if err != nil {
		log.Fatal(err)
	}

	bullet, _, err = ebitenutil.NewImageFromFile("image/Bullet1.png")
	if err != nil {
		log.Fatal(err)
	}

	// フォントの定義
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	fontFace, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

// 1/60 秒で呼ばれる関数
func (g *Game) Update() error {
	// モードによって遷移
	switch {
	// メニュー
	case mode == 0:
		if g.KeyPressed_Space() {
			mode = 1
		}
	// アニメーション
	case mode == 1:
		for i := 0; i < 50; i++ {
			knightx = knightx + 0.01
			ghostx = ghostx - 0.01
		}
		if knightx >= 70 {
			random = rand.Float64()*5 + 4
			now = time.Now()
			mode = 2
		}
	// Ready の時間の処理
	case mode == 2:
		stop := time.Now()
		elapsedTime = stop.Sub(now).Seconds()
		// お手付き
		if g.KeyPressed_Space() {
			mode = 5
		}
		if elapsedTime > random {
			now = time.Now()
			mode = 3
		}
	// Start の時間の処理
	case mode == 3:
		if g.KeyPressed_Space() {
			stop := time.Now()
			sub = stop.Sub(now).Seconds()
			// 勝ち
			if ghostTime >= sub {
				judgement = true
				knightBool = false
			} else {
				// 負け
				judgement = false
			}
			mode = 4
		}
	// 勝ち処理
	case (mode == 4) && judgement:
		score = score + 100
		time.Sleep(3 * time.Second)
		ghostTime = ghostTime / 2
		knightx = 0
		ghostx = 280
		knightBool = true
		sub = 0
		mode = 1
	// 負け処理
	case (mode == 4) && !judgement:
		if g.KeyPressed_Space() {
			if highScore < score {
				Save()
			} else {
				os.Exit(0)
			}
		}
	// お手付き処理
	case mode == 5:
		time.Sleep(3 * time.Second)
		otetuki++
		// お手付き２回でGameOver
		if otetuki == 2 {
			judgement = false
			mode = 4
		} else {
			knightx = 0
			ghostx = 280
			mode = 1
		}
	}
	return nil
}

// 描画に関する処理
func (g *Game) Draw(screen *ebiten.Image) {
	text.Draw(screen, strconv.FormatFloat(sub, 'f', -1, 64), fontFace, 100, 200, color.RGBA{255, 255, 255, 255}) // ストップした時間の表示

	// 騎士の描写
	if knightBool {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(knightx, 150)
		op.GeoM.Scale(2, 2)
		screen.DrawImage(knight, op)
	}

	// ゴーストの描写
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ghostx, 150)
	op.GeoM.Scale(2, 2)
	screen.DrawImage(ghost, op)

	// モードで描写遷移
	switch {
	// メニュー
	case mode == 0:
		text.Draw(screen, "SamuraiGo", fontFace, 230, 200, color.RGBA{255, 255, 255, 255})
		text.Draw(screen, "SPACE : START", fontFace, 200, 350, color.RGBA{255, 0, 0, 255})
	// Ready
	case mode == 2:
		text.Draw(screen, "Ready", fontFace, 250, 100, color.RGBA{255, 255, 255, 255})
	// Start
	case mode == 3:
		text.Draw(screen, "START!!!!", fontFace, 250, 100, color.RGBA{255, 0, 0, 255})
	// Win
	case mode == 4 && judgement:
		text.Draw(screen, "Win", fontFace, 250, 100, color.RGBA{255, 0, 0, 255})

		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(180, 140)
		op.GeoM.Scale(2, 2)
		screen.DrawImage(bullet, op)

		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(260, 150)
		op.GeoM.Scale(2, 2)
		screen.DrawImage(knightWin, op)
	// GAME OVER
	case mode == 4 && !judgement:
		screen.Fill(color.Black)
		text.Draw(screen, "GAME OVER", fontFace, 230, 200, color.RGBA{255, 0, 0, 255})
		text.Draw(screen, fmt.Sprintf("score: %d", score), fontFace, 230, 250, color.White)
		text.Draw(screen, "SPACE : SAVE & EXIT", fontFace, 150, 350, color.White)
	// お手付き
	case mode == 5:
		text.Draw(screen, "Too Early", fontFace, 230, 200, color.RGBA{255, 0, 0, 255})
	}
}

// スペースキーの判定
func (g *Game) KeyPressed_Space() bool {
	return inpututil.IsKeyJustPressed(103)
}

// スクリーンのレイアウト
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

// ハイスコア読み込み
func Load() {
	file, err := os.Open("text/score.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	for s.Scan() {
		ss := strings.Split(s.Text(), ",")
		highScore, err = strconv.Atoi(ss[0])

		if err != nil {
			fmt.Fprintln(os.Stderr, "エラー：", err)
			os.Exit(1)
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println("HighScore : ", highScore)
}

// セーブ
func Save() {
	file, err := os.Create("text/score.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(strconv.Itoa(score))
	writer.Flush()
	os.Exit(0)
}

func main() {
	Load()
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("SamuraiGo")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
