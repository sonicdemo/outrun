package netobj

import (
	"log"
	"math/rand"
	"strconv"

	"github.com/Mtbcooler/outrun/config"
	"github.com/Mtbcooler/outrun/consts"
	"github.com/Mtbcooler/outrun/enums"
	"github.com/Mtbcooler/outrun/obj"
	"github.com/jinzhu/now"
)

type WheelOptions struct {
	Items                []string   `json:"items"`
	Item                 []int64    `json:"item"`
	ItemWeight           []int64    `json:"itemWeight"`
	ItemWon              int64      `json:"itemWon"`
	NextFreeSpin         int64      `json:"nextFreeSpin"` // midnight (start of next day)
	SpinCost             int64      `json:"spinCost"`
	RouletteRank         int64      `json:"rouletteRank"`
	NumRouletteToken     int64      `json:"numRouletteToken"`
	NumJackpotRing       int64      `json:"numJackpotRing"`
	NumRemainingRoulette int64      `json:"numRemainingRoulette"`
	ItemList             []obj.Item `json:"itemList"`
}

func DefaultWheelOptions(numRouletteTicket, rouletteCountInPeriod, rouletteRank int64) WheelOptions {
	// TODO: Modifying this seems like a good way of figuring out what the game thinks each ID means in terms of items.
	// NOTE: Free spins occur when numRemainingRoulette > numRouletteToken
	// TODO: Add different randomly-picked "modes" for roulette generation
	items := []string{strconv.Itoa(enums.IDTypeItemRouletteWin)} // first item is always jackpot/big/super
	itemAmounts := []int64{1}
	for _ = range make([]byte, 7) { // loop 7 times
		randomItem := consts.RandomItemListNormalWheel[rand.Intn(len(consts.RandomItemListNormalWheel))]
		randomItemAmount := consts.NormalWheelItemAmountRange[randomItem].GetRandom()
		items = append(items, randomItem)
		itemAmounts = append(itemAmounts, randomItemAmount)
	}
	/*
	   itemAmounts := []int64{1, 2, 2, 2, 1, 3, 2, 2}
	*/
	itemWeight := []int64{1250, 1250, 1250, 1250, 1250, 1250, 1250, 1250}
	//itemWon := int64(0)
	itemWon := int64(rand.Intn(len(items)))
	nextFreeSpin := now.EndOfDay().Unix() + 1 // midnight
	spinCost := int64(87)
	//rouletteRank := int64(enums.WheelRankNormal)
	//numRouletteToken := playerState.NumRouletteTicket
	numRouletteToken := numRouletteTicket                // The game uses the _current_ value, not as if it was in the past (This is hard to explain, maybe TODO: explain this better?)
	numJackpotRing := int64(consts.RouletteJackpotRings) // TODO: Make jackpot value dynamic
	// TODO: get rid of logic here!
	numRemainingRoulette := numRouletteToken + consts.RouletteFreeSpins - rouletteCountInPeriod // TODO: is this proper?
	if numRemainingRoulette < numRouletteToken {
		numRemainingRoulette = numRouletteToken
	}
	itemList := []obj.Item{}
	out := WheelOptions{
		items,
		itemAmounts,
		itemWeight,
		itemWon,
		nextFreeSpin,
		spinCost,
		rouletteRank,
		numRouletteToken,
		numJackpotRing,
		numRemainingRoulette,
		itemList,
	}
	return out
}

func UpgradeWheelOptions(origWheel WheelOptions, numRouletteTicket, rouletteCountInPeriod int64) WheelOptions {
	rouletteRank := origWheel.RouletteRank
	if origWheel.Items[origWheel.ItemWon] == strconv.Itoa(enums.IDTypeItemRouletteWin) { // if landed on big/super or jackpot
		landedOnUpgrade := origWheel.RouletteRank == enums.WheelRankNormal || origWheel.RouletteRank == enums.WheelRankBig
		if config.CFile.DebugPrints {
			log.Printf("%v\n", origWheel.RouletteRank)
			log.Printf("%v\n", landedOnUpgrade)
		}
		if landedOnUpgrade {
			if config.CFile.DebugPrints {
				log.Println("landedOnUpgrade")
			}
			rouletteRank++ // increase the rank
		} else {
			if config.CFile.DebugPrints {
				log.Println("NOT landedOnUpgrade")
			}
			rouletteRank = enums.WheelRankNormal
		}
	}
	newWheel := DefaultWheelOptions(numRouletteTicket, rouletteCountInPeriod, rouletteRank)
	return newWheel
}
