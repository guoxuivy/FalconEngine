package plugins

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"utils"
)

type ButTimesPlus struct {
	BuyRules Buyrules
}

type Order struct {

	Count    int64  `json:"count"`
	TotalAmount   float64  `json:"total_amount"`
	RealAmount	  float64  `json:"real_amount"`	
}


type Buyrules struct {
	StartDate string `json:"start"`
	EndDate   string `json:"end"`
	Count     int64  `json:"count"`
	Opration  string `json:"op"`
}

func NewBuyTimes() *ButTimesPlus {
	var BuyRules Buyrules
	this := &ButTimesPlus{BuyRules: BuyRules}
	return this
}


func (this *ButTimesPlus) Init() bool {
	return true
}


func (this *ButTimesPlus) SetRules(rules interface{}) func(value_byte interface{}) bool {
	fmt.Printf("rules : %v \n",rules)
	rule,ok:=rules.(utils.Condition)
	if !ok {
		fmt.Printf("Error rules\n")
	}
	var start,end string
	date_range := strings.Split(rule.Range,",")
	if len(date_range) != 2{
		start= "2015-01-01"
		end = "2015-12-31"
	}
	start = date_range[0]
	end = date_range[1]
	var total_count int64
	var total_amount float64
	var err error
	if rule.Key == "buy_count"{
		total_count, err = strconv.ParseInt(rule.Value, 0, 0)
		if err != nil {
			fmt.Printf("Error %v \n", rule.Value)
		}
	}else{
		total_amount, err = strconv.ParseFloat(rule.Value,0)
		if err != nil {
			fmt.Printf("Error %v \n", rule.Value)
		}		
	}
	
	//fmt.Printf("total : %v start : %v end : %v \n",total,start,end)
	return func(value_byte interface{}) bool{
		var err error
		buytimes := make(map[string]Order)
		body, ok := value_byte.([]byte)
		if !ok {
			fmt.Printf("Byte Error ...\n")
		}
		err = json.Unmarshal(body, &buytimes)
		if err != nil {
			fmt.Printf("Unmarshal Error ...\n")
			return false
		}
		
		var sum int64 = 0
		var sum_amount float64 = 0.0
		for date,value := range buytimes{
			//fmt.Printf("date : %v start : %v end : %v sum : %v count : %v \n",date,start,end,sum,value.Count)
			if date > start  && date < end {
				sum = sum + value.Count
				sum_amount = sum_amount + value.RealAmount
			}
		}
		switch rule.Op{
			case "more":
				return ((sum > total_count && rule.Key == "buy_count") || (sum_amount > total_amount && rule.Key == "buy_amount"))
			case "less":
				return ((sum < total_count && rule.Key == "buy_count") || (sum_amount < total_amount && rule.Key == "buy_amount"))
			case "equal":
				return ((sum == total_count && rule.Key == "buy_count") || (sum_amount == total_amount && rule.Key == "buy_amount"))
		}
		//if (sum > total_count && rule.Key == "buy_count") || (sum_amount > total_amount && rule.Key == "buy_amount"){
			//fmt.Printf("Match .... %v \n", buytimes)
		//	return true
		//}
		//fmt.Printf("Not Match .... \n")
		return false
	}
}



func (this *ButTimesPlus) CustomeFunction(v1, v2 interface{}) bool {
	/*
	var err error
	var buytimes Buytimes
	body, ok := v2.([]byte)
	if !ok {
		fmt.Printf("Byte Error ...\n")
	}
	err = json.Unmarshal(body, &buytimes)
	if err != nil {
		fmt.Printf("Unmarshal Error ...\n")
		return false
	}
	var sum int64 = 0
	for i, _ := range buytimes.BuyDetail {
		if buytimes.BuyDetail[i].DateTime > "2015-03-05" {
			sum = sum + buytimes.BuyDetail[i].Count
		}
	}
	if sum > 5 {
		//fmt.Printf("Match .... %v \n", buytimes)
		return true
	}
	//fmt.Printf("Not Match .... \n")
	*/
	return false
}



//插件分词函数,返回string数组,bool参数表示是建立索引的时候还是查询的调用,STYPE = 9 调用
func (this *ButTimesPlus) SegmentFunc(value interface{},isSearch bool) []string{
	res := make([]string,0)
	if isSearch == true{
		res=append(res,fmt.Sprintf("%v",value))
		return res
	}
	
	
	//fmt.Printf("SegmentFunc...\n")
	buytimes := make(map[string]Order)
	body, ok := value.(string)
	if !ok {
		fmt.Printf("Byte Error ...\n")
	}
	err := json.Unmarshal([]byte(body), &buytimes)
	if err != nil {
		fmt.Printf("Unmarshal Error ...\n")
		return nil
	}
	
	for date,_ := range buytimes{
		//fmt.Printf("date : %v  value  : %v \n",date,value)
		res=append(res,date)
		
	} 
	//fmt.Printf("res : %v \n",res)
	return res
}


//数字分词函数,返回string数组,bool参数表示是建立索引的时候还是查询的调用,STYPE = 9 调用
func (this *ButTimesPlus) SplitNum(value interface{}) int64{
	
	return 0
}




//插件正排处理函数，建立索引的时候调用，stype =9 调用 ,返回byte数组
func (this *ButTimesPlus) BuildByteProfile(value []byte) []byte {
	
	return value
}

//插件正排处理函数，建立索引的时候调用，stype =9 调用 ,返回string,定长！！！！
func (this *ButTimesPlus) BuildStringProfile(value interface{}) string{
	
	return "nil"
} 



//插件正排处理函数，建立索引的时候调用，stype =9 调用 ,返回int64
func (this *ButTimesPlus) BuildIntProfile(value interface{}) int64{
	
	return 0
}