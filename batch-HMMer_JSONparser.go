package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type hitstruct []struct {
	Flags    int    `json:"flags"`
	Nregions int    `json:"nregions"`
	Ndom     int    `json:"ndom"`
	Name     string `json:"name"`
	Score    string `json:"score"`
	Bias     string `json:"bias"`
	Taxid    string `json:"taxid"`
	Acc      string `json:"acc"`
	Domains  []struct {
		Alisqacc    string  `json:"alisqacc"`
		AliIDCount  int     `json:"aliIdCount"`
		Alirfline   string  `json:"alirfline"`
		IsIncluded  int     `json:"is_included"`
		Alihmmname  string  `json:"alihmmname"`
		Bitscore    float64 `json:"bitscore"`
		Display     int     `json:"display"`
		Ievalue     string  `json:"ievalue"`
		Alisqto     int     `json:"alisqto"`
		AliSim      float64 `json:"aliSim"`
		Jali        int     `json:"jali"`
		Bias        string  `json:"bias"`
		Ienv        int     `json:"ienv"`
		Cevalue     string  `json:"cevalue"`
		Significant int     `json:"significant"`
		Alimline    string  `json:"alimline"`
		Alihmmfrom  int     `json:"alihmmfrom"`
		Clan        string  `json:"clan"`
		AliL        int     `json:"aliL"`
		Alihindex   string  `json:"alihindex"`
		IsReported  int     `json:"is_reported"`
		Alintseq    string  `json:"alintseq"`
		Jenv        int     `json:"jenv"`
		Alimmline   string  `json:"alimmline"`
		Alihmmacc   string  `json:"alihmmacc"`
		Oasc        string  `json:"oasc"`
		Aliaseq     string  `json:"aliaseq"`
		Alihmmto    int     `json:"alihmmto"`
		AliID       float64 `json:"aliId"`
		Alippline   string  `json:"alippline"`
		Alimodel    string  `json:"alimodel"`
		AliM        int     `json:"aliM"`
		Iali        int     `json:"iali"`
		Alicsline   string  `json:"alicsline"`
		AliSimCount int     `json:"aliSimCount"`
		Alihmmdesc  string  `json:"alihmmdesc"`
		Alisqdesc   string  `json:"alisqdesc"`
		Outcompeted int     `json:"outcompeted"`
		Alisqname   string  `json:"alisqname"`
		Alisqfrom   int     `json:"alisqfrom"`
		Uniq        int     `json:"uniq"`
		AliN        int     `json:"aliN"`
	} `json:"domains"`
	Nincluded int     `json:"nincluded"`
	Evalue    string  `json:"evalue"`
	Desc      string  `json:"desc"`
	Pvalue    float64 `json:"pvalue"`
	Nreported int     `json:"nreported"`
	Hindex    string  `json:"hindex"`
}

type hmmerResult struct {
	Results struct {
		Hits  hitstruct `json:"hits"`
		Stats struct {
			Nhits     int    `json:"nhits"`
			Elapsed   string `json:"elapsed"`
			Z         int    `json:"Z"`
			ZSetby    int    `json:"Z_setby"`
			NPastMsv  int    `json:"n_past_msv"`
			Unpacked  int    `json:"unpacked"`
			User      int    `json:"user"`
			DomZSetby int    `json:"domZ_setby"`
			Nseqs     int    `json:"nseqs"`
			NPastBias int    `json:"n_past_bias"`
			Sys       int    `json:"sys"`
			NPastFwd  int    `json:"n_past_fwd"`
			Nmodels   int    `json:"nmodels"`
			Nincluded int    `json:"nincluded"`
			NPastVit  int    `json:"n_past_vit"`
			Nreported int    `json:"nreported"`
			DomZ      int    `json:"domZ"`
		} `json:"stats"`
	} `json:"results"`
}

func (h hmmerResult) String() string {
	formater := "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v"
	if h.Results.Stats.Nhits == 0 {
		return "empty"
	}
	if len(h.Results.Hits) == 0 {
		return "no hits, but lists >0 in stats.nhits\tweird"
	}
	if h.Results.Hits[0].Domains[0].Clan == "" {
		h.Results.Hits[0].Domains[0].Clan = "No_CLAN"
	}
	return fmt.Sprintf(formater,
		strings.ReplaceAll(h.Results.Hits[0].Domains[0].Alisqname, ">", ""),
		h.Results.Stats.Nhits,
		h.Results.Hits[0].Name,
		h.Results.Hits[0].Acc,
		h.Results.Hits[0].Domains[0].Clan,
		h.Results.Hits[0].Domains[0].Ienv,
		h.Results.Hits[0].Domains[0].Jenv,
		h.Results.Hits[0].Domains[0].Iali,
		h.Results.Hits[0].Domains[0].Jali,
		h.Results.Hits[0].Domains[0].Alihmmfrom,
		h.Results.Hits[0].Domains[0].Alihmmto,
		h.Results.Hits[0].Domains[0].Bitscore,
		h.Results.Hits[0].Domains[0].Ievalue,
		h.Results.Hits[0].Domains[0].Cevalue,
		h.Results.Hits[0].Desc)
}

func main() {
	var location string
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println(location + "wazaaa")
	// 	}
	// }()
	var receiveBucket hmmerResult
	title := "file\tnmatch\tQuery Name\tNumber of Hits\tFamily id\tFamily Accession\tClan\tEnv. Start\tEnv. End\tAli. Start\tAli. End\tModel Start\tModel End\tBit Score\tInd. E-value\tCond. E-value\tDescription"
	fmt.Println(title)
	matches, err := WalkMatch(".", "*json")
	if err != nil {
		panic(err)
	}
	var nmatch int
	for nmatch, location = range matches {
		file, err := ioutil.ReadFile(location)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(file, &receiveBucket)
		fmt.Printf("%v\t%v\t%v\n", location, nmatch+1, receiveBucket)
	}
}

//The function below will recursively walk through a directory and return the paths to all files whose name matches the given pattern:

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
