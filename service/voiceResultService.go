package service

import (
	"encoding/xml"
	"log"
)

type xmlResult struct {
	XMLName     xml.Name     `xml:"xml_result"`
	ReadChapter readChapter0 `xml:"read_chapter"`
}

type readChapter0 struct {
	XMLName  xml.Name `xml:"read_chapter"`
	Lan      string   `xml:"lan,attr"`
	Type     string   `xml:"type,attr"`
	Version  string   `xml:"version,attr"`
	RecPaper recPaper `xml:"rec_paper"`
}

type recPaper struct {
	XMLName     xml.Name    `xml:"rec_paper"`
	ReadChapter ReadChapter `xml:"read_chapter"`
}

// ReadChapter 篇章信息
type ReadChapter struct {
	XMLName        xml.Name   `xml:"read_chapter"`
	AccuracyScore  float32    `xml:"accuracy_score,attr"`
	BegPos         int        `xml:"beg_pos,attr"`
	Content        string     `xml:"content,attr"`
	EngPos         int        `xml:"eng_pos,attr"`
	ExceptInfo     int        `xml:"except_info,attr"`
	FluencyScore   float32    `xml:"fluency_score,attr"`
	IntegrityScore float32    `xml:"integrity_score,attr"`
	IsRejected     bool       `xml:"is_rejected,attr"`
	RejectType     int        `xml:"reject_type,attr"`
	ScorePattern   string     `xml:"score_pattern,attr"`
	StandardScore  float32    `xml:"standard_score,attr"`
	TotalScore     float32    `xml:"total_score,attr"`
	WordCount      int        `xml:"word_count,attr"`
	Sentence       []sentence `xml:"sentence"`
}

// 句子信息
type sentence struct {
	XMLName       xml.Name `xml:"sentence"`
	AccuracyScore float32  `xml:"accuracy_score,attr"`
	BegPos        int      `xml:"beg_pos,attr"`
	Content       string   `xml:"content,attr"`
	EndPos        int      `xml:"end_pos,attr"`
	FluencyScore  float32  `xml:"fluency_score,attr"`
	Index         int      `xml:"index,attr"`
	StandardScore float32  `xml:"standard_score,attr"`
	TotalScore    float32  `xml:"total_score,attr"`
	WordCount     int      `xml:"word_count,attr"`
	Word          []word   `xml:"word"`
}

// 单词信息
type word struct {
	XMLName     xml.Name `xml:"word"`
	BegPos      int      `xml:"beg_pos,attr"`
	Content     string   `xml:"content,attr"`
	DpMessage   int      `xml:"dp_message,attr"`
	EndPos      int      `xml:"end_pos,attr"`
	GlobalIndex int      `xml:"global_index,attr"`
	Index       int      `xml:"index,attr"`
	Property    int      `xml:"property,attr"`
	TotalScore  float32  `xml:"total_score,attr"`
	Syll        []syll   `xml:"syll"`
}

// 音节
type syll struct {
	XMLName     xml.Name `xml:"syll"`
	BegPos      int      `xml:"beg_pos,attr"`
	Content     string   `xml:"content,attr"`
	EndPos      int      `xml:"end_pos,attr"`
	SerrMsg     int      `xml:"serr_msg,attr"`
	SyllAccount int      `xml:"syll_account,attr"`
	SyllScore   float32  `xml:"syll_score,attr"`
	Phone       []phone  `xml:"phone"`
}

// 音素
type phone struct {
	XMLName xml.Name `xml:"phone"`
	BegPos  int      `xml:"beg_pos,attr"`
	Content string   `xml:"content,attr"`
	EndPos  int      `xml:"end_pos,attr"`
}

// HandleVoiceXML 解析返回的xml信息
func HandleVoiceXML(voiceXML []byte) ReadChapter {
	XMLResult := xmlResult{}
	err := xml.Unmarshal(voiceXML, &XMLResult)
	if err != nil {
		log.Fatal(err)
	}
	result := XMLResult.ReadChapter.RecPaper.ReadChapter
	for index1 := 0; index1 < len(result.Sentence); index1++ {
		target := &result.Sentence[index1].Word
		for index2 := 0; index2 < len(*target); index2++ {
			if (*target)[index2].Content == "sil" {
				*target = append((*target)[:index2], (*target)[index2+1:]...)
				index2--
			}
		}
	}
	return result
}
