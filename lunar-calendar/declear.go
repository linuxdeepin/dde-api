/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

type lunarYearInfo struct {
	leapMonth     int //闰月所在月，0为没有闰月
	zhengMonth    int //正月初一对应公历月
	zhengDay      int //正月初一对应公历日
	lunarMonthNum int //农历每月的天数的数组(需转换为二进制,得到每月大小，0=小月(29日),1=大月(30日))
}

type caDayInfo struct {
	index int
	days  int
}

type caYearInfo struct {
	Year  int
	Month int
	Day   int
}

type caLunarDayInfo struct {
	LunarYear      int
	LunarMonth     int
	LunarDay       int
	LunarLeapMonth int
	LunarMonthName string
	LunarDayName   string
	GanZhiYear     string
	GanZhiMonth    string
	GanZhiDay      string
	Zodiac         string
	Term           string
	SolarFestival  string
	LunarFestival  string
	Worktime       int
}

type caSolarMonthInfo struct {
	FirstDayWeek int
	Days         int
	Datas        []caYearInfo
}

type caLunarMonthInfo struct {
	FirstDayWeek int
	Days         int
	Datas        []caLunarDayInfo
}

type cacheUtil struct {
	current int
}

var (
	MinYear  = 1890 //最小年限
	MaxYear  = 2100 //最大年限
	cacheMap = make(map[string]interface{})
	cacheObj = newCache()

	lunarData = map[string][]string{
		"heavenlyStems":   {"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"},                                                                                               //天干
		"earthlyBranches": {"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"},                                                                                     //地支
		"zodiac":          {"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"},                                                                                     //对应地支十二生肖
		"solarTerm":       {"小寒", "大寒", "立春", "雨水", "惊蛰", "春分", "清明", "谷雨", "立夏", "小满", "芒种", "夏至", "小暑", "大暑", "立秋", "处暑", "白露", "秋分", "寒露", "霜降", "立冬", "小雪", "大雪", "冬至"}, //二十四节气
		"monthCn":         {"正", "二", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "十二"},
		"dateCn":          {"初一", "初二", "初三", "初四", "初五", "初六", "初七", "初八", "初九", "初十", "十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九", "二十", "廿一", "廿二", "廿三", "廿四", "廿五", "廿六", "廿七", "廿八", "廿九", "三十", "卅一"},
	}

	//公历节日
	solarFestival = map[string]string{
		"d0101": "元旦节",
		"d0202": "世界湿地日",
		"d0210": "国际气象节",
		"d0214": "情人节",
		"d0301": "国际海豹日",
		"d0303": "全国爱耳日",
		"d0305": "学雷锋纪念日",
		"d0308": "妇女节",
		"d0312": "植树节 孙中山逝世纪念日",
		"d0314": "国际警察日",
		"d0315": "消费者权益日",
		"d0317": "中国国医节 国际航海日",
		"d0321": "世界森林日 消除种族歧视国际日 世界儿歌日",
		"d0322": "世界水日",
		"d0323": "世界气象日",
		"d0324": "世界防治结核病日",
		"d0325": "全国中小学生安全教育日",
		"d0330": "巴勒斯坦国土日",
		"d0401": "愚人节 全国爱国卫生运动月(四月) 税收宣传月(四月)",
		"d0407": "世界卫生日",
		"d0422": "世界地球日",
		"d0423": "世界图书和版权日",
		"d0424": "亚非新闻工作者日",
		"d0501": "劳动节",
		"d0504": "青年节",
		"d0505": "碘缺乏病防治日",
		"d0508": "世界红十字日",
		"d0512": "国际护士节",
		"d0515": "国际家庭日",
		"d0517": "世界电信日",
		"d0518": "国际博物馆日",
		"d0520": "全国学生营养日",
		"d0522": "国际生物多样性日",
		"d0523": "国际牛奶日",
		"d0531": "世界无烟日",
		"d0601": "国际儿童节",
		"d0605": "世界环境日",
		"d0606": "全国爱眼日",
		"d0617": "防治荒漠化和干旱日",
		"d0623": "国际奥林匹克日",
		"d0625": "全国土地日",
		"d0626": "国际禁毒日",
		"d0701": "香港回归纪念日 中共诞辰 世界建筑日",
		"d0702": "国际体育记者日",
		"d0707": "抗日战争纪念日",
		"d0711": "世界人口日",
		"d0730": "非洲妇女日",
		"d0801": "建军节",
		"d0808": "中国男子节(爸爸节)",
		"d0815": "抗日战争胜利纪念",
		"d0908": "国际扫盲日 国际新闻工作者日",
		"d0909": "毛泽东逝世纪念",
		"d0910": "中国教师节",
		"d0914": "世界清洁地球日",
		"d0916": "国际臭氧层保护日",
		"d0918": "九一八事变纪念日",
		"d0920": "国际爱牙日",
		"d0927": "世界旅游日",
		"d0928": "孔子诞辰",
		"d1001": "国庆节 世界音乐日 国际老人节",
		"d1002": "国际和平与民主自由斗争日",
		"d1004": "世界动物日",
		"d1006": "老人节",
		"d1008": "全国高血压日 世界视觉日",
		"d1009": "世界邮政日 万国邮联日",
		"d1010": "辛亥革命纪念日 世界精神卫生日",
		"d1013": "世界保健日 国际教师节",
		"d1014": "世界标准日",
		"d1015": "国际盲人节(白手杖节)",
		"d1016": "世界粮食日",
		"d1017": "世界消除贫困日",
		"d1022": "世界传统医药日",
		"d1024": "联合国日 世界发展信息日",
		"d1031": "世界勤俭日",
		"d1107": "十月社会主义革命纪念日",
		"d1108": "中国记者日",
		"d1109": "全国消防安全宣传教育日",
		"d1110": "世界青年节",
		"d1111": "国际科学与和平周(本日所属的一周)",
		"d1112": "孙中山诞辰纪念日",
		"d1114": "世界糖尿病日",
		"d1117": "国际大学生节 世界学生节",
		"d1121": "世界问候日 世界电视日",
		"d1129": "国际声援巴勒斯坦人民国际日",
		"d1201": "世界艾滋病日",
		"d1203": "世界残疾人日",
		"d1205": "国际经济和社会发展志愿人员日",
		"d1208": "国际儿童电视日",
		"d1209": "世界足球日",
		"d1210": "世界人权日",
		"d1212": "西安事变纪念日",
		"d1213": "南京大屠杀(1937年)纪念日！紧记血泪史！",
		"d1220": "澳门回归纪念",
		"d1221": "国际篮球日",
		"d1224": "平安夜",
		"d1225": "圣诞节",
		"d1226": "毛泽东诞辰纪念",
	}

	//农历节日
	lunarFestival = map[string]string{
		"d0101": "春节",
		"d0115": "元宵节",
		"d0202": "龙抬头节",
		"d0323": "妈祖生辰",
		"d0505": "端午节",
		"d0707": "七夕情人节",
		"d0715": "中元节",
		"d0815": "中秋节",
		"d0909": "重阳节",
		"d1015": "下元节",
		"d1208": "腊八节",
		"d1223": "小年",
		"d0100": "除夕",
	}

	lunarInfos = []lunarYearInfo{
		lunarYearInfo{2, 1, 21, 22184},
		lunarYearInfo{0, 2, 9, 21936},
		lunarYearInfo{6, 1, 30, 9656},
		lunarYearInfo{0, 2, 17, 9584},
		lunarYearInfo{0, 2, 6, 21168},
		lunarYearInfo{5, 1, 26, 43344},
		lunarYearInfo{0, 2, 13, 59728},
		lunarYearInfo{0, 2, 2, 27296},
		lunarYearInfo{3, 1, 22, 44368},
		lunarYearInfo{0, 2, 10, 43856},
		lunarYearInfo{8, 1, 30, 19304},
		lunarYearInfo{0, 2, 19, 19168},
		lunarYearInfo{0, 2, 8, 42352},
		lunarYearInfo{5, 1, 29, 21096},
		lunarYearInfo{0, 2, 16, 53856},
		lunarYearInfo{0, 2, 4, 55632},
		lunarYearInfo{4, 1, 25, 27304},
		lunarYearInfo{0, 2, 13, 22176},
		lunarYearInfo{0, 2, 2, 39632},
		lunarYearInfo{2, 1, 22, 19176},
		lunarYearInfo{0, 2, 10, 19168},
		lunarYearInfo{6, 1, 30, 42200},
		lunarYearInfo{0, 2, 18, 42192},
		lunarYearInfo{0, 2, 6, 53840},
		lunarYearInfo{5, 1, 26, 54568},
		lunarYearInfo{0, 2, 14, 46400},
		lunarYearInfo{0, 2, 3, 54944},
		lunarYearInfo{2, 1, 23, 38608},
		lunarYearInfo{0, 2, 11, 38320},
		lunarYearInfo{7, 2, 1, 18872},
		lunarYearInfo{0, 2, 20, 18800},
		lunarYearInfo{0, 2, 8, 42160},
		lunarYearInfo{5, 1, 28, 45656},
		lunarYearInfo{0, 2, 16, 27216},
		lunarYearInfo{0, 2, 5, 27968},
		lunarYearInfo{4, 1, 24, 44456},
		lunarYearInfo{0, 2, 13, 11104},
		lunarYearInfo{0, 2, 2, 38256},
		lunarYearInfo{2, 1, 23, 18808},
		lunarYearInfo{0, 2, 10, 18800},
		lunarYearInfo{6, 1, 30, 25776},
		lunarYearInfo{0, 2, 17, 54432},
		lunarYearInfo{0, 2, 6, 59984},
		lunarYearInfo{5, 1, 26, 27976},
		lunarYearInfo{0, 2, 14, 23248},
		lunarYearInfo{0, 2, 4, 11104},
		lunarYearInfo{3, 1, 24, 37744},
		lunarYearInfo{0, 2, 11, 37600},
		lunarYearInfo{7, 1, 31, 51560},
		lunarYearInfo{0, 2, 19, 51536},
		lunarYearInfo{0, 2, 8, 54432},
		lunarYearInfo{6, 1, 27, 55888},
		lunarYearInfo{0, 2, 15, 46416},
		lunarYearInfo{0, 2, 5, 22176},
		lunarYearInfo{4, 1, 25, 43736},
		lunarYearInfo{0, 2, 13, 9680},
		lunarYearInfo{0, 2, 2, 37584},
		lunarYearInfo{2, 1, 22, 51544},
		lunarYearInfo{0, 2, 10, 43344},
		lunarYearInfo{7, 1, 29, 46248},
		lunarYearInfo{0, 2, 17, 27808},
		lunarYearInfo{0, 2, 6, 46416},
		lunarYearInfo{5, 1, 27, 21928},
		lunarYearInfo{0, 2, 14, 19872},
		lunarYearInfo{0, 2, 3, 42416},
		lunarYearInfo{3, 1, 24, 21176},
		lunarYearInfo{0, 2, 12, 21168},
		lunarYearInfo{8, 1, 31, 43344},
		lunarYearInfo{0, 2, 18, 59728},
		lunarYearInfo{0, 2, 8, 27296},
		lunarYearInfo{6, 1, 28, 44368},
		lunarYearInfo{0, 2, 15, 43856},
		lunarYearInfo{0, 2, 5, 19296},
		lunarYearInfo{4, 1, 25, 42352},
		lunarYearInfo{0, 2, 13, 42352},
		lunarYearInfo{0, 2, 2, 21088},
		lunarYearInfo{3, 1, 21, 59696},
		lunarYearInfo{0, 2, 9, 55632},
		lunarYearInfo{7, 1, 30, 23208},
		lunarYearInfo{0, 2, 17, 22176},
		lunarYearInfo{0, 2, 6, 38608},
		lunarYearInfo{5, 1, 27, 19176},
		lunarYearInfo{0, 2, 15, 19152},
		lunarYearInfo{0, 2, 3, 42192},
		lunarYearInfo{4, 1, 23, 53864},
		lunarYearInfo{0, 2, 11, 53840},
		lunarYearInfo{8, 1, 31, 54568},
		lunarYearInfo{0, 2, 18, 46400},
		lunarYearInfo{0, 2, 7, 46752},
		lunarYearInfo{6, 1, 28, 38608},
		lunarYearInfo{0, 2, 16, 38320},
		lunarYearInfo{0, 2, 5, 18864},
		lunarYearInfo{4, 1, 25, 42168},
		lunarYearInfo{0, 2, 13, 42160},
		lunarYearInfo{10, 2, 2, 45656},
		lunarYearInfo{0, 2, 20, 27216},
		lunarYearInfo{0, 2, 9, 27968},
		lunarYearInfo{6, 1, 29, 44448},
		lunarYearInfo{0, 2, 17, 43872},
		lunarYearInfo{0, 2, 6, 38256},
		lunarYearInfo{5, 1, 27, 18808},
		lunarYearInfo{0, 2, 15, 18800},
		lunarYearInfo{0, 2, 4, 25776},
		lunarYearInfo{3, 1, 23, 27216},
		lunarYearInfo{0, 2, 10, 59984},
		lunarYearInfo{8, 1, 31, 27432},
		lunarYearInfo{0, 2, 19, 23232},
		lunarYearInfo{0, 2, 7, 43872},
		lunarYearInfo{5, 1, 28, 37736},
		lunarYearInfo{0, 2, 16, 37600},
		lunarYearInfo{0, 2, 5, 51552},
		lunarYearInfo{4, 1, 24, 54440},
		lunarYearInfo{0, 2, 12, 54432},
		lunarYearInfo{0, 2, 1, 55888},
		lunarYearInfo{2, 1, 22, 23208},
		lunarYearInfo{0, 2, 9, 22176},
		lunarYearInfo{7, 1, 29, 43736},
		lunarYearInfo{0, 2, 18, 9680},
		lunarYearInfo{0, 2, 7, 37584},
		lunarYearInfo{5, 1, 26, 51544},
		lunarYearInfo{0, 2, 14, 43344},
		lunarYearInfo{0, 2, 3, 46240},
		lunarYearInfo{4, 1, 23, 46416},
		lunarYearInfo{0, 2, 10, 44368},
		lunarYearInfo{9, 1, 31, 21928},
		lunarYearInfo{0, 2, 19, 19360},
		lunarYearInfo{0, 2, 8, 42416},
		lunarYearInfo{6, 1, 28, 21176},
		lunarYearInfo{0, 2, 16, 21168},
		lunarYearInfo{0, 2, 5, 43312},
		lunarYearInfo{4, 1, 25, 29864},
		lunarYearInfo{0, 2, 12, 27296},
		lunarYearInfo{0, 2, 1, 44368},
		lunarYearInfo{2, 1, 22, 19880},
		lunarYearInfo{0, 2, 10, 19296},
		lunarYearInfo{6, 1, 29, 42352},
		lunarYearInfo{0, 2, 17, 42208},
		lunarYearInfo{0, 2, 6, 53856},
		lunarYearInfo{5, 1, 26, 59696},
		lunarYearInfo{0, 2, 13, 54576},
		lunarYearInfo{0, 2, 3, 23200},
		lunarYearInfo{3, 1, 23, 27472},
		lunarYearInfo{0, 2, 11, 38608},
		lunarYearInfo{11, 1, 31, 19176},
		lunarYearInfo{0, 2, 19, 19152},
		lunarYearInfo{0, 2, 8, 42192},
		lunarYearInfo{6, 1, 28, 53848},
		lunarYearInfo{0, 2, 15, 53840},
		lunarYearInfo{0, 2, 4, 54560},
		lunarYearInfo{5, 1, 24, 55968},
		lunarYearInfo{0, 2, 12, 46496},
		lunarYearInfo{0, 2, 1, 22224},
		lunarYearInfo{2, 1, 22, 19160},
		lunarYearInfo{0, 2, 10, 18864},
		lunarYearInfo{7, 1, 30, 42168},
		lunarYearInfo{0, 2, 17, 42160},
		lunarYearInfo{0, 2, 6, 43600},
		lunarYearInfo{5, 1, 26, 46376},
		lunarYearInfo{0, 2, 14, 27936},
		lunarYearInfo{0, 2, 2, 44448},
		lunarYearInfo{3, 1, 23, 21936},
		lunarYearInfo{0, 2, 11, 37744},
		lunarYearInfo{8, 2, 1, 18808},
		lunarYearInfo{0, 2, 19, 18800},
		lunarYearInfo{0, 2, 8, 25776},
		lunarYearInfo{6, 1, 28, 27216},
		lunarYearInfo{0, 2, 15, 59984},
		lunarYearInfo{0, 2, 4, 27424},
		lunarYearInfo{4, 1, 24, 43872},
		lunarYearInfo{0, 2, 12, 43744},
		lunarYearInfo{0, 2, 2, 37600},
		lunarYearInfo{3, 1, 21, 51568},
		lunarYearInfo{0, 2, 9, 51552},
		lunarYearInfo{7, 1, 29, 54440},
		lunarYearInfo{0, 2, 17, 54432},
		lunarYearInfo{0, 2, 5, 55888},
		lunarYearInfo{5, 1, 26, 23208},
		lunarYearInfo{0, 2, 14, 22176},
		lunarYearInfo{0, 2, 3, 42704},
		lunarYearInfo{4, 1, 23, 21224},
		lunarYearInfo{0, 2, 11, 21200},
		lunarYearInfo{8, 1, 31, 43352},
		lunarYearInfo{0, 2, 19, 43344},
		lunarYearInfo{0, 2, 7, 46240},
		lunarYearInfo{6, 1, 27, 46416},
		lunarYearInfo{0, 2, 15, 44368},
		lunarYearInfo{0, 2, 5, 21920},
		lunarYearInfo{4, 1, 24, 42448},
		lunarYearInfo{0, 2, 12, 42416},
		lunarYearInfo{0, 2, 2, 21168},
		lunarYearInfo{3, 1, 22, 43320},
		lunarYearInfo{0, 2, 9, 26928},
		lunarYearInfo{7, 1, 29, 29336},
		lunarYearInfo{0, 2, 17, 27296},
		lunarYearInfo{0, 2, 6, 44368},
		lunarYearInfo{5, 1, 26, 19880},
		lunarYearInfo{0, 2, 14, 19296},
		lunarYearInfo{0, 2, 3, 42352},
		lunarYearInfo{4, 1, 24, 21104},
		lunarYearInfo{0, 2, 10, 53856},
		lunarYearInfo{8, 1, 30, 59696},
		lunarYearInfo{0, 2, 18, 54560},
		lunarYearInfo{0, 2, 7, 55968},
		lunarYearInfo{6, 1, 27, 27472},
		lunarYearInfo{0, 2, 15, 22224},
		lunarYearInfo{0, 2, 5, 19168},
		lunarYearInfo{4, 1, 25, 42216},
		lunarYearInfo{0, 2, 12, 42192},
		lunarYearInfo{0, 2, 1, 53584},
		lunarYearInfo{2, 1, 21, 55592},
		lunarYearInfo{0, 2, 9, 54560},
	}

	/**
	 * 二十四节气数据，节气点时间（单位是分钟）
	 * 从0小寒起算
	 */
	termInfo = []int{0, 21208, 42467, 63836, 85337, 107014, 128867, 150921, 173149, 195551, 218072, 240693, 263343, 285989, 308563, 331033, 353350, 375494, 397447, 419210, 440795, 462224, 483532, 504758}

	//中国节日放假安排，外部设置，0无特殊安排，1工作，2放假
	worktimeYearMap = map[string]map[string]int{
		"y2013": map[string]int{
			"d0101": 2,
			"d0102": 2,
			"d0103": 2,
			"d0105": 1,
			"d0106": 1,
			"d0209": 2,
			"d0210": 2,
			"d0211": 2,
			"d0212": 2,
			"d0213": 2,
			"d0214": 2,
			"d0215": 2,
			"d0216": 1,
			"d0217": 1,
			"d0404": 2,
			"d0405": 2,
			"d0406": 2,
			"d0407": 1,
			"d0427": 1,
			"d0428": 1,
			"d0429": 2,
			"d0430": 2,
			"d0501": 2,
			"d0608": 1,
			"d0609": 1,
			"d0610": 2,
			"d0611": 2,
			"d0612": 2,
			"d0919": 2,
			"d0920": 2,
			"d0921": 2,
			"d0922": 1,
			"d0929": 1,
			"d1001": 2,
			"d1002": 2,
			"d1003": 2,
			"d1004": 2,
			"d1005": 2,
			"d1006": 2,
			"d1007": 2,
			"d1012": 1,
		},
		"y2014": map[string]int{
			"d0101": 2,
			"d0126": 1,
			"d0131": 2,
			"d0201": 2,
			"d0202": 2,
			"d0203": 2,
			"d0204": 2,
			"d0205": 2,
			"d0206": 2,
			"d0208": 1,
			"d0405": 2,
			"d0407": 2,
			"d0501": 2,
			"d0502": 2,
			"d0503": 2,
			"d0504": 1,
			"d0602": 2,
			"d0908": 2,
			"d0928": 1,
			"d1001": 2,
			"d1002": 2,
			"d1003": 2,
			"d1004": 2,
			"d1005": 2,
			"d1006": 2,
			"d1007": 2,
			"d1011": 1,
		},
	}
)
