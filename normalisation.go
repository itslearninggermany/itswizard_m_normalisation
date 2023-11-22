package itswizard_m_normalisation

import (
	"fmt"
	itswizard_basic "github.com/itslearninggermany/itswizard_m_basic"
	bw_structs "github.com/itslearninggermany/itswizard_m_bwStructs"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
	"unicode"
)

type specialCharacters struct {
	old        string
	runeNumber string
	number     int
}

func checkData(buf []byte) (sc []specialCharacters) {
	for i, rune := range string(buf) {

		if strings.Contains(strconv.QuoteRuneToASCII(rune), "\\u") {
			sc = append(sc, specialCharacters{
				runeNumber: fmt.Sprint(rune),
				old:        fmt.Sprint(string(rune)),
				number:     i,
			})
		}
	}
	return sc
}

func Normalise(data []byte) []byte {
	sc := checkData(data)

	str := string(data)

	for i := 0; i < len(sc); i++ {
		str = strings.Replace(str, sc[i].old, Characters[sc[i].runeNumber], sc[i].number)
	}
	return []byte(str)
}

func NormaliseForEmail(data []byte) []byte {
	sc := checkData(data)

	str := string(data)

	for i := 0; i < len(sc); i++ {
		str = strings.Replace(str, sc[i].old, CharactersEmail[sc[i].runeNumber], sc[i].number)
	}

	return []byte(str)
}

var CharactersEmail = map[string]string{
	"223": "ss", // ß
	"246": "oe", // ö
	"228": "ae", // ä
	"252": "ue", // ü

	"233": "e", // é

	"253":  "y", // ý 253
	"7845": "a", // ấ 7845
	"7879": "e", // ệ 7879
}

var Characters = map[string]string{
	"253":  "y", // ý 253
	"7845": "a", // ấ 7845
	"7879": "e", // ệ 7879
}

func Utf8ToUsername(input string) string {
	out := ""
	for _, rune := range input {
		nu := fmt.Sprint(CharactersToChange[int(rune)])
		if nu != "" {
			out = out + nu
		} else {
			out = out + fmt.Sprint(string(rune))
		}
	}
	return out
}

func Utf8ForFirstAndLastname(input string) string {
	out := ""
	for _, rune := range input {
		nu := fmt.Sprint(CharactersToChangeIntWithUmlaut[int(rune)])
		if nu != "" {
			out = out + nu
		} else {
			out = out + fmt.Sprint(int(rune))
		}
	}
	return out
}

func CreateUsernameForLusd(firstname, lastname string, db *gorm.DB) string {
	username := ""
	for i, k := range firstname {
		if i > 2 {
			break
		}
		username = username + string(k)
	}

	username = username + "."
	for i, k := range lastname {
		if i > 2 {
			break
		}
		username = username + string(k)
	}

	i := 0

	var tmpusername = username
	for {
		var tmp itswizard_basic.LusdPerson
		if i == 0 {
			if db.Where("username = ?", tmpusername).Last(&tmp).RecordNotFound() {
				break
			} else {
				i++
				continue
			}
		}

		tmpusername = username + strconv.Itoa(i)
		db.Where("username = ?", tmpusername).Last(&tmp)

		if tmp.Username == "" {
			username = tmpusername
			break
		}

		i++
	}

	return strings.ToLower(Utf8ToUsername(username))
}

func CreateUsernameForBw(firstname, lastname string, db *gorm.DB, allNewUsernames map[string]bool) string {
	username := ""
	for i, k := range SpaceMap(firstname) {
		if i > 3 {
			break
		}
		username = username + string(k)
	}

	username = username + "."
	for i, k := range SpaceMap(lastname) {
		if i > 3 {
			break
		}
		username = username + string(k)
	}

	i := 0
	username = strings.ToLower(Utf8ForFirstAndLastname(username))

	var tmpusername = username
	for {
		var tmp bw_structs.BWPerson
		if i == 0 {
			if db.Where("username = ?", tmpusername).Last(&tmp).RecordNotFound() {
				for {
					if allNewUsernames[tmpusername] == false {
						return tmpusername
					} else {
						i++
						tmpusername = username + strconv.Itoa(i)
					}
				}
			} else {
				i++
				continue
			}
		}
		tmpusername = username + strconv.Itoa(i)
		db.Where("username = ?", tmpusername).Last(&tmp)
		if tmp.Username == "" {
			for {
				if allNewUsernames[tmpusername] == false {
					return tmpusername
				} else {
					i++
					tmpusername = username + strconv.Itoa(i)
				}
			}
		}
		i++
	}
	return "wrongUsername"
}

func SpaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

var CharactersToChangeIntWithUmlaut = map[int]string{
	0:    "0",     //NUL	0	0000	null character
	1:    "1",     //SOH	1	0001	start of header
	2:    "2",     //STX	2	0002	start of text
	3:    "3",     //ETX	3	0003	end of text
	4:    "4",     //EOT	4	0004	end of transmission
	5:    "",      //ENQ	5	0005	enquiry
	6:    "",      //ACK	6	0006	acknowledge
	7:    "",      //BEL	7	0007	bell (ring)
	8:    "",      //BS	8	0008	backspace
	9:    "",      //HT	9	0009	horizontal tab
	10:   "",      //LF	10	000A	line feed
	11:   "",      //VT	11	000B	vertical tab
	12:   "",      //FF	12	000C	form feed
	13:   "",      //CR	13	000D	carriage return
	14:   "",      //SO	14	000E	shift out
	15:   "",      //SI	15	000F	shift in
	16:   "",      //DLE	16	0010	data link escape
	17:   "",      //DC1	17	0011	device control 1
	18:   "",      //DC2	18	0012	device control 2
	19:   "",      //DC3	19	0013	device control 3
	20:   "",      //DC4	20	0014	device control 4
	21:   "",      //NAK	21	0015	negative acknowledge
	22:   "",      //SYN	22	0016	synchronize
	23:   "",      //ETB	23	0017	end transmission block
	24:   "",      //CAN	24	0018	cancel
	25:   "",      //EM	25	0019	end of medium
	26:   "",      //SUB	26	001A	substitute
	27:   "",      //ESC	27	001B	escape
	28:   "",      //FS	28	001C	file separator
	29:   "",      //GS	29	001D	group separator
	30:   "",      //RS	30	001E	record separator
	31:   "",      //US	31	001F	unit separator
	32:   " ",     //	32	0020	 	SPACE
	33:   "",      //!	33	0021	 	EXCLAMATION MARK
	34:   "",      //"	34	0022	&quot;	QUOTATION MARK
	35:   "",      //#	35	0023	 	NUMBER SIGN
	36:   "",      //$	36	0024	 	DOLLAR SIGN
	37:   "",      //%	37	0025	 	PERCENT SIGN
	38:   "",      //&	38	0026	&amp;	AMPERSAND
	39:   "",      //'	39	0027	 	APOSTROPHE
	40:   "",      //(	40	0028	 	LEFT PARENTHESIS
	41:   "",      //)	41	0029	 	RIGHT PARENTHESIS
	42:   "",      //*	42	002A	 	ASTERISK
	43:   "",      //+	43	002B	 	PLUS SIGN
	44:   "",      //,	44	002C	 	COMMA
	45:   "-",     //-	45	002D	 	HYPHEN-MINUS
	46:   ".",     //.	46	002E	 	FULL STOP
	47:   "",      ///	47	002F	 	SOLIDUS
	48:   "0",     //0	48	0030	 	DIGIT ZERO
	49:   "1",     //1	49	0031	 	DIGIT ONE
	50:   "2",     //2	50	0032	 	DIGIT TWO
	51:   "3",     //3 0033	 	DIGIT THREE
	52:   "4",     //4 0034	 	DIGIT FOUR
	53:   "5",     //5 0035	 	DIGIT FIVE
	54:   "6",     //6 0036	 	DIGIT SIX
	55:   "7",     //7 0037	 	DIGIT SEVEN
	56:   "8",     //8 0038	 	DIGIT EIGHT
	57:   "9",     //9 0039	 	DIGIT NINE
	58:   "",      //:	003A	 	COLON
	59:   "",      //;	003B	 	SEMICOLON
	60:   "",      //<	003C	&lt;	LESS-THAN SIGN
	61:   "",      //=	003D	 	EQUALS SIGN
	62:   "",      //>	003E	&gt;	GREATER-THAN SIGN
	63:   "",      //?	003F	 	QUESTION MARK
	64:   "",      //@	0040	 	COMMERCIAL AT
	65:   "A",     //A	0041	 	LATIN CAPITAL LETTER A
	66:   "B",     //B	0042	 	LATIN CAPITAL LETTER B
	67:   "C",     //C	0043	 	LATIN CAPITAL LETTER C
	68:   "D",     //D	0044	 	LATIN CAPITAL LETTER D
	69:   "E",     //E 0045	 	LATIN CAPITAL LETTER E
	70:   "F",     //F	70	0046	 	LATIN CAPITAL LETTER F
	71:   "G",     //G	71	0047	 	LATIN CAPITAL LETTER G
	72:   "H",     //H	72	0048	 	LATIN CAPITAL LETTER H
	73:   "I",     //I	73	0049	 	LATIN CAPITAL LETTER I
	74:   "J",     //J	74	004A	 	LATIN CAPITAL LETTER J
	75:   "K",     //K	75	004B	 	LATIN CAPITAL LETTER K
	76:   "L",     //L 76	004C	 	LATIN CAPITAL LETTER L
	77:   "M",     //M	77	004D	 	LATIN CAPITAL LETTER M
	78:   "N",     //N	78	004E	 	LATIN CAPITAL LETTER N
	79:   "O",     //O	79	004F	 	LATIN CAPITAL LETTER O
	80:   "P",     //P	80	0050	 	LATIN CAPITAL LETTER P
	81:   "Q",     //Q	81	0051	 	LATIN CAPITAL LETTER Q
	82:   "R",     //R	82	0052	 	LATIN CAPITAL LETTER R
	83:   "S",     //S	83	0053	 	LATIN CAPITAL LETTER S
	84:   "T",     //T	84	0054	 	LATIN CAPITAL LETTER T
	85:   "U",     //U	85	0055	 	LATIN CAPITAL LETTER U
	86:   "V",     //V	86	0056	 	LATIN CAPITAL LETTER V
	87:   "W",     //W	87	0057	 	LATIN CAPITAL LETTER W
	88:   "X",     //X	88	0058	 	LATIN CAPITAL LETTER X
	89:   "Y",     //Y	89	0059	 	LATIN CAPITAL LETTER Y
	90:   "Z",     //Z	90	005A	 	LATIN CAPITAL LETTER Z
	91:   "",      // [	91	005B	 	LEFT SQUARE BRACKET
	92:   "",      // \	92	005C	 	REVERSE SOLIDUS
	93:   "",      //]	93	005D	 	RIGHT SQUARE BRACKET
	94:   "",      // ^	94	005E	 	CIRCUMFLEX ACCENT
	95:   "",      //_	95	005F	 	LOW LINE
	96:   "",      //`	96	0060	 	GRAVE ACCENT
	97:   "a",     // a	97	0061	 	LATIN SMALL LETTER A
	98:   "b",     // b	98	0062	 	LATIN SMALL LETTER B
	99:   "c",     //c	99	0063	 	LATIN SMALL LETTER C
	100:  "d",     //d	100	0064	 	LATIN SMALL LETTER D
	101:  "e",     //e	101	0065	 	LATIN SMALL LETTER E
	102:  "f",     //f	102	0066	 	LATIN SMALL LETTER F
	103:  "g",     //g	103	0067	 	LATIN SMALL LETTER G
	104:  "h",     //h	104	0068	 	LATIN SMALL LETTER H
	105:  "i",     //i	105	0069	 	LATIN SMALL LETTER I
	106:  "j",     //j	106	006A	 	LATIN SMALL LETTER J
	107:  "k",     //k	107	006B	 	LATIN SMALL LETTER K
	108:  "l",     //l	108	006C	 	LATIN SMALL LETTER L
	109:  "m",     //m	109	006D	 	LATIN SMALL LETTER M
	110:  "n",     //n	110	006E	 	LATIN SMALL LETTER N
	111:  "o",     //o	111	006F	 	LATIN SMALL LETTER O
	112:  "p",     //p	112	0070	 	LATIN SMALL LETTER P
	113:  "q",     //q	113	0071	 	LATIN SMALL LETTER Q
	114:  "r",     //r	114	0072	 	LATIN SMALL LETTER R
	115:  "s",     //s	115	0073	 	LATIN SMALL LETTER S
	116:  "t",     //t	116	0074	 	LATIN SMALL LETTER T
	117:  "u",     // u	117	0075	 	LATIN SMALL LETTER U
	118:  "v",     //v	118	0076	 	LATIN SMALL LETTER V
	119:  "w",     //w	119	0077	 	LATIN SMALL LETTER W
	120:  "x",     //x	120	0078	 	LATIN SMALL LETTER X
	121:  "y",     //y	121	0079	 	LATIN SMALL LETTER Y
	122:  "z",     //z	122	007A	 	LATIN SMALL LETTER Z
	123:  "",      //{	123	007B	 	LEFT CURLY BRACKET
	124:  "",      //|	124	007C	 	VERTICAL LINE
	125:  "",      //}	125	007D	 	RIGHT CURLY BRACKET
	126:  "",      //~	126	007E	 	TILDE
	127:  "",      // DEL	127	007F	delete (rubout)
	128:  "",      //€	128	0080	CONTROL
	129:  "",      // 	129	0081	CONTROL
	130:  "",      //‚	130	0082	BREAK PERMITTED HERE
	131:  "",      //ƒ	131	0083	NO BREAK HERE
	132:  "",      //„	132	0084	INDEX
	133:  "",      //…	133	0085	NEXT LINE (NEL)
	134:  "",      //†	134	0086	START OF SELECTED AREA
	135:  "",      //‡	135	0087	END OF SELECTED AREA
	136:  "",      //ˆ	136	0088	CHARACTER TABULATION SET
	137:  "",      //‰	137	0089	CHARACTER TABULATION WITH JUSTIFICATION
	138:  "",      //Š	138	008A	LINE TABULATION SET
	139:  "",      //‹	139	008B	PARTIAL LINE FORWARD
	140:  "",      //Œ	140	008C	PARTIAL LINE BACKWARD
	141:  "",      //	141	008D	REVERSE LINE FEED
	142:  "",      //Ž	142	008E	SINGLE SHIFT TWO
	143:  "",      //	143	008F	SINGLE SHIFT THREE
	144:  "",      //	144	0090	DEVICE CONTROL STRING
	145:  "",      //‘	145	0091	PRIVATE USE ONE
	146:  "",      //’	146	0092	PRIVATE USE TWO
	147:  "",      //“	147	0093	SET TRANSMIT STATE
	148:  "",      //”	148	0094	CANCEL CHARACTER
	149:  "",      //•	149	0095	MESSAGE WAITING
	150:  "",      //–	150	0096	START OF GUARDED AREA
	151:  "",      //—	151	0097	END OF GUARDED AREA
	152:  "",      //˜	152	0098	START OF STRING
	153:  "",      //™	153	0099	CONTROL
	154:  "",      //š	154	009A	SINGLE CHARACTER INTRODUCER
	155:  "",      //›	155	009B	CONTROL SEQUENCE INTRODUCER
	156:  "",      //œ	156	009C	STRING TERMINATOR
	157:  "",      //	157	009D	OPERATING SYSTEM COMMAND
	158:  "",      //ž	158	009E	PRIVACY MESSAGE
	159:  "",      //Ÿ	159	009F	APPLICATION PROGRAM COMMAND
	160:  " ",     //	160	00A0	&nbsp;	NO-BREAK SPACE
	161:  "",      //¡	161	00A1	&iexcl;	INVERTED EXCLAMATION MARK
	162:  "",      //¢	162	00A2	&cent;	CENT SIGN
	163:  "",      //£	163	00A3	&pound;	POUND SIGN
	164:  "",      //¤	164	00A4	&curren;	CURRENCY SIGN
	165:  "",      //¥	165	00A5	&yen;	YEN SIGN
	166:  "",      //¦	166	00A6	&brvbar;	BROKEN BAR
	167:  "",      //§	167	00A7	&sect;	SECTION SIGN
	168:  "",      //¨	168	00A8	&uml;	DIAERESIS
	169:  "",      //©	169	00A9	&copy;	COPYRIGHT SIGN
	170:  "",      //ª	170	00AA	&ordf;	FEMININE ORDINAL INDICATOR
	171:  "",      //«	171	00AB	&laquo;	LEFT-POINTING DOUBLE ANGLE QUOTATION MARK
	172:  "",      //¬	172	00AC	&not;	NOT SIGN
	173:  "",      //­	173	00AD	&shy;	SOFT HYPHEN
	174:  "",      //®	174	00AE	&reg;	REGISTERED SIGN
	175:  "",      //¯	175	00AF	&macr;	MACRON
	176:  "",      //°	176	00B0	&deg;	DEGREE SIGN
	177:  "",      //±	177	00B1	&plusmn;	PLUS-MINUS SIGN
	178:  "",      //²	178	00B2	&sup2;	SUPERSCRIPT TWO
	179:  "",      //³	179	00B3	&sup3;	SUPERSCRIPT THREE
	180:  "",      //´	180	00B4	&acute;	ACUTE ACCENT
	181:  "",      //µ	181	00B5	&micro;	MICRO SIGN
	182:  "",      //¶	182	00B6	&para;	PILCROW SIGN
	183:  "",      //·	183	00B7	&middot;	MIDDLE DOT
	184:  "",      //¸	184	00B8	&cedil;	CEDILLA
	185:  "",      //¹	185	00B9	&sup1;	SUPERSCRIPT ONE
	186:  "",      //º	186	00BA	&ordm;	MASCULINE ORDINAL INDICATOR
	187:  "",      //»	187	00BB	&raquo;	RIGHT-POINTING DOUBLE ANGLE QUOTATION MARK
	188:  "",      //¼	188	00BC	&frac14;	VULGAR FRACTION ONE QUARTER
	189:  "",      //½	189	00BD	&frac12;	VULGAR FRACTION ONE HALF
	190:  "",      //¾	190	00BE	&frac34;	VULGAR FRACTION THREE QUARTERS
	191:  "",      //¿	191	00BF	&iquest;	INVERTED QUESTION MARK
	192:  "A",     // À	192	00C0	&Agrave;	LATIN CAPITAL LETTER A WITH GRAVE
	193:  "A",     // Á	193
	194:  "A",     // Â	194
	195:  "A",     // Ã	195
	196:  "A",     // Ä	196
	197:  "A",     // Å	197
	198:  "A",     // Æ	198
	199:  "C",     //Ç	199	00C7	&Ccedil;	LATIN CAPITAL LETTER C WITH CEDILLA
	200:  "E",     //È	200	00C8	&Egrave;	LATIN CAPITAL LETTER E WITH GRAVE
	201:  "E",     //É	201	00C9	&Eacute;	LATIN CAPITAL LETTER E WITH ACUTE
	202:  "E",     //Ê	202	00CA	&Ecirc;	LATIN CAPITAL LETTER E WITH CIRCUMFLEX
	203:  "E",     //Ë	203	00CB	&Euml;	LATIN CAPITAL LETTER E WITH DIAERESIS
	204:  "I",     //Ì	204	00CC	&Igrave;	LATIN CAPITAL LETTER I WITH GRAVE
	205:  "I",     //Í	205	00CD	&Iacute;	LATIN CAPITAL LETTER I WITH ACUTE
	206:  "I",     //Î	206	00CE	&Icirc;	LATIN CAPITAL LETTER I WITH CIRCUMFLEX
	207:  "I",     //Ï	207	00CF	&Iuml;	LATIN CAPITAL LETTER I WITH DIAERESIS
	208:  "D",     //Ð	208	00D0	&ETH;	LATIN CAPITAL LETTER ETH
	209:  "N",     //Ñ	209	00D1	&Ntilde;	LATIN CAPITAL LETTER N WITH TILDE
	210:  "O",     //Ò	210	00D2	&Ograve;	LATIN CAPITAL LETTER O WITH GRAVE
	211:  "O",     //Ó	211	00D3	&Oacute;	LATIN CAPITAL LETTER O WITH ACUTE
	212:  "O",     //Ô	212	00D4	&Ocirc;	LATIN CAPITAL LETTER O WITH CIRCUMFLEX
	213:  "O",     //Õ	213	00D5	&Otilde;	LATIN CAPITAL LETTER O WITH TILDE
	214:  "Oe",    //Ö	214	00D6	&Ouml;	LATIN CAPITAL LETTER O WITH DIAERESIS
	215:  "",      //×	215	00D7	&times;	MULTIPLICATION SIGN
	216:  "O",     //Ø	216	00D8	&Oslash;	LATIN CAPITAL LETTER O WITH STROKE
	217:  "U",     //Ù	217	00D9	&Ugrave;	LATIN CAPITAL LETTER U WITH GRAVE
	218:  "U",     //Ú	218	00DA	&Uacute;	LATIN CAPITAL LETTER U WITH ACUTE
	219:  "U",     //Û	219	00DB	&Ucirc;	LATIN CAPITAL LETTER U WITH CIRCUMFLEX
	220:  "U",     //Ü	220	00DC	&Uuml;	LATIN CAPITAL LETTER U WITH DIAERESIS
	221:  "Y",     // Ý	221	00DD	&Yacute;	LATIN CAPITAL LETTER Y WITH ACUTE
	222:  "",      // Þ	222	00DE	&THORN;	LATIN CAPITAL LETTER THORN
	223:  "ss",    // ß
	224:  "a",     //à	224	00E0	&agrave;	LATIN SMALL LETTER A WITH GRAVE
	225:  "a",     //á	225	00E1	&aacute;	LATIN SMALL LETTER A WITH ACUTE
	226:  "a",     //â	226	00E2	&acirc;	LATIN SMALL LETTER A WITH CIRCUMFLEX
	227:  "a",     //ã	227	00E3	&atilde;	LATIN SMALL LETTER A WITH TILDE
	228:  "ae",    // ä
	229:  "a",     // å	229	00E5	&aring;	LATIN SMALL LETTER A WITH RING ABOVE
	230:  "ae",    // æ	230	00E6	&aelig;	LATIN SMALL LETTER AE
	231:  "c",     //ç	231	00E7	&ccedil;	LATIN SMALL LETTER C WITH CEDILLA
	232:  "e",     // è
	233:  "e",     // é
	234:  "e",     // ê
	235:  "e",     // 	ë	235	00EB	&euml;	LATIN SMALL LETTER E WITH DIAERESIS
	236:  "i",     //ì	236	00EC	&igrave;	LATIN SMALL LETTER I WITH GRAVE
	237:  "i",     //í	237	00ED	&iacute;	LATIN SMALL LETTER I WITH ACUTE
	238:  "i",     //î	238	00EE	&icirc;	LATIN SMALL LETTER I WITH CIRCUMFLEX
	239:  "i",     //ï	239	00EF	&iuml;	LATIN SMALL LETTER I WITH DIAERESIS
	240:  "o",     //ï	239	00EF	&iuml;	LATIN SMALL LETTER I WITH DIAERESIS
	241:  "n",     //	ñ	241	00F1	&ntilde;	LATIN SMALL LETTER N WITH TILDE
	242:  "o",     //ò	242	00F2	&ograve;	LATIN SMALL LETTER O WITH GRAVE
	243:  "o",     //ó	243	00F3	&oacute;	LATIN SMALL LETTER O WITH ACUTE
	244:  "o",     //ô	244	00F4	&ocirc;	LATIN SMALL LETTER O WITH CIRCUMFLEX
	245:  "o",     //õ	245	00F5	&otilde;	LATIN SMALL LETTER O WITH TILDE
	246:  "oe",    // ö
	247:  "",      //÷	247	00F7	&divide;	DIVISION SIGN
	248:  "o",     //ø	248	00F8	&oslash;	LATIN SMALL LETTER O WITH STROKE
	249:  "u",     //ù	249	00F9	&ugrave;	LATIN SMALL LETTER U WITH GRAVE
	250:  "u",     //ú	250	00FA	&uacute;	LATIN SMALL LETTER U WITH ACUTE
	251:  "u",     //û	251	00FB	&ucirc;	LATIN SMALL LETTER U WITH CIRCUMFLEX
	252:  "ue",    // ü
	253:  "y",     // ý 253
	254:  "b",     //þ	254	00FE	&thorn;	LATIN SMALL LETTER THORN
	255:  "y",     //ÿ	255	00FF	&yuml;	LATIN SMALL LETTER Y WITH DIAERESIS
	256:  "A",     //Ā	256	0100	&Amacr;	LATIN CAPITAL LETTER A WITH MACRON
	257:  "a",     //ā	257	0101	&amacr;	LATIN SMALL LETTER A WITH MACRON
	258:  "A",     //Ă	258	0102	&Abreve;	LATIN CAPITAL LETTER A WITH BREVE
	259:  "a",     //ă	259	0103	&abreve;	LATIN SMALL LETTER A WITH BREVE
	260:  "A",     //Ą	260	0104	&Aogon;	LATIN CAPITAL LETTER A WITH OGONEK
	261:  "q",     //q	261	0105	&aogon;	LATIN SMALL LETTER A WITH OGONEK
	262:  "C",     //Ć	262	0106	&Cacute;	LATIN CAPITAL LETTER C WITH ACUTE
	263:  "c",     //ć	263	0107	&cacute;	LATIN SMALL LETTER C WITH ACUTE
	264:  "C",     //Ĉ	264	0108	&Ccirc;	LATIN CAPITAL LETTER C WITH CIRCUMFLEX
	265:  "c",     //ĉ	265	0109	&ccirc;	LATIN SMALL LETTER C WITH CIRCUMFLEX
	266:  "C",     //Ċ	266	010A	&Cdod;	LATIN CAPITAL LETTER C WITH DOT ABOVE
	267:  "c",     //ċ	267	010B	&cdot;	LATIN SMALL LETTER C WITH DOT ABOVE
	268:  "C",     //Č	268	010C	&Ccaron;	LATIN CAPITAL LETTER C WITH CARON
	269:  "c",     //č	269	010D	&ccaron;	LATIN SMALL LETTER C WITH CARON
	270:  "C",     //Ď	270	010E	&Dcaron;	LATIN CAPITAL LETTER D WITH CARON
	271:  "c",     //ď	271	010F	&dcaron;	LATIN SMALL LETTER D WITH CARON
	272:  "C",     //Đ	272	0110	&Dstrok;	LATIN CAPITAL LETTER D WITH STROKE
	273:  "c",     //đ	273	0111	&dstrok;	LATIN SMALL LETTER D WITH STROKE
	274:  "E",     //Ē	274	0112	&Emacr;	LATIN CAPITAL LETTER E WITH MACRON
	275:  "e",     //ē	275	0113	&emacr;	LATIN SMALL LETTER E WITH MACRON
	276:  "E",     //Ĕ	276	0114	 	LATIN CAPITAL LETTER E WITH BREVE
	277:  "e",     //ĕ	277	0115	 	LATIN SMALL LETTER E WITH BREVE
	278:  "E",     //Ė	278	0116	&Edot;	LATIN CAPITAL LETTER E WITH DOT ABOVE
	279:  "e",     //ė	279	0117	&edot;	LATIN SMALL LETTER E WITH DOT ABOVE
	280:  "E",     //Ę	280	0118	&Eogon;	LATIN CAPITAL LETTER E WITH OGONEK
	281:  "e",     //ę	281	0119	&eogon;	LATIN SMALL LETTER E WITH OGONEK
	282:  "E",     //Ě	282	011A	&Ecaron;	LATIN CAPITAL LETTER E WITH CARON
	283:  "e",     //ě	283	011B	&ecaron;	LATIN SMALL LETTER E WITH CARON
	284:  "G",     //Ĝ	284	011C	&Gcirc;	LATIN CAPITAL LETTER G WITH CIRCUMFLEX
	285:  "g",     //ĝ	285	011D	&gcirc;	LATIN SMALL LETTER G WITH CIRCUMFLEX
	286:  "G",     //Ğ	286	011E	&Gbreve;	LATIN CAPITAL LETTER G WITH BREVE
	287:  "g",     //ğ	287	011F	&gbreve;	LATIN SMALL LETTER G WITH BREVE
	288:  "G",     //Ġ	288	0120	&Gdot;	LATIN CAPITAL LETTER G WITH DOT ABOVE
	289:  "g",     //ġ	289	0121	&gdot;	LATIN SMALL LETTER G WITH DOT ABOVE
	290:  "G",     //Ģ	290	0122	&Gcedil;	LATIN CAPITAL LETTER G WITH CEDILLA
	291:  "g",     //ģ	291	0123	&gcedil;	LATIN SMALL LETTER G WITH CEDILLA
	292:  "H",     //Ĥ	292	0124	&Hcirc;	LATIN CAPITAL LETTER H WITH CIRCUMFLEX
	293:  "h",     //ĥ	293	0125	&hcirc;	LATIN SMALL LETTER H WITH CIRCUMFLEX
	294:  "H",     //Ħ	294	0126	&Hstrok;	LATIN CAPITAL LETTER H WITH STROKE
	295:  "h",     //ħ	295	0127	&hstrok;	LATIN SMALL LETTER H WITH STROKE
	296:  "i",     // Ĩ &Itilde;	LATIN CAPITAL LETTER I WITH TILDE
	297:  "i",     // ĩ	297	0129	&itilde;	LATIN SMALL LETTER I WITH TILDE
	298:  "i",     // Ī	298	012A	&Imacr;	LATIN CAPITAL LETTER I WITH MACRON
	299:  "i",     // ī	299	012B	&imacr;	LATIN SMALL LETTER I WITH MACRON
	300:  "i",     // Ĭ	300	012C	 	LATIN CAPITAL LETTER I WITH BREVE
	301:  "i",     // ĭ	301	012D	 	LATIN SMALL LETTER I WITH BREVE
	302:  "i",     // Į	302	012E	&Iogon;	LATIN CAPITAL LETTER I WITH OGONEK
	303:  "i",     // į	303	012F	&iogon;	LATIN SMALL LETTER I WITH OGONEK
	304:  "i",     // İ	304	0130	&Idot;	LATIN CAPITAL LETTER I WITH DOT ABOVE
	305:  "i",     // ı	305	0131	&inodot;	LATIN SMALL LETTER DOTLESS I
	306:  "ij",    // Ĳ	306	0132	&IJlog;	LATIN CAPITAL LIGATURE IJ
	307:  "ij",    // ĳ	307	0133	&ijlig;	LATIN SMALL LIGATURE IJ
	308:  "j",     //Ĵ	308	0134	&Jcirc;	LATIN CAPITAL LETTER J WITH CIRCUMFLEX
	309:  "j",     //ĵ	309	0135	&jcirc;	LATIN SMALL LETTER J WITH CIRCUMFLEX
	310:  "K",     // Ķ	310	0136	&Kcedil;	LATIN CAPITAL LETTER K WITH CEDILLA
	311:  "k",     //ķ	311	0137	&kcedli;	LATIN SMALL LETTER K WITH CEDILLA
	312:  "k",     //ĸ	312	0138	&kgreen;	LATIN SMALL LETTER KRA
	313:  "L",     // 	Ĺ	313	0139	&Lacute;	LATIN CAPITAL LETTER L WITH ACUTE
	314:  "l",     // 	ĺ	314	013A	&lacute;	LATIN SMALL LETTER L WITH ACUTE
	315:  "L",     //Ļ	315	013B	&Lcedil;	LATIN CAPITAL LETTER L WITH CEDILLA
	316:  "l",     //ļ	316	013C	&lcedil;	LATIN SMALL LETTER L WITH CEDILLA
	317:  "L",     //Ľ	317	013D	&Lcaron;	LATIN CAPITAL LETTER L WITH CARON
	318:  "l",     //ľ	318	013E	&lcaron;	LATIN SMALL LETTER L WITH CARON
	319:  "L",     //Ŀ	319	013F	&Lmodot;	LATIN CAPITAL LETTER L WITH MIDDLE DOT
	320:  "l",     //ŀ	320	0140	&lmidot;	LATIN SMALL LETTER L WITH MIDDLE DOT
	321:  "L",     //Ł	321	0141	&Lstrok;	LATIN CAPITAL LETTER L WITH STROKE
	322:  "l",     //ł	322	0142	&lstrok;	LATIN SMALL LETTER L WITH STROKE
	323:  "N",     //Ń	323	0143	&Nacute;	LATIN CAPITAL LETTER N WITH ACUTE
	324:  "n",     //ń	324	0144	&nacute;	LATIN SMALL LETTER N WITH ACUTE
	325:  "N",     //Ņ	325	0145	&Ncedil;	LATIN CAPITAL LETTER N WITH CEDILLA
	326:  "n",     //ņ	326	0146	&ncedil;	LATIN SMALL LETTER N WITH CEDILLA
	327:  "N",     //Ň	327	0147	&Ncaron;	LATIN CAPITAL LETTER N WITH CARON
	328:  "n",     //ň	328	0148	&ncaron;	LATIN SMALL LETTER N WITH CARON
	329:  "n",     //ŉ	329	0149	&napos;	LATIN SMALL LETTER N PRECEDED BY APOSTROPHE
	330:  "N",     //Ŋ	330	014A	&ENG;	LATIN CAPITAL LETTER ENG
	331:  "n",     //ŋ	331	014B	&eng;	LATIN SMALL LETTER ENG
	332:  "O",     //Ō	332	014C	&Omacr;	LATIN CAPITAL LETTER O WITH MACRON
	333:  "o",     //ō	333	014D	&omacr;	LATIN SMALL LETTER O WITH MACRON
	334:  "O",     //Ŏ	334	014E	 	LATIN CAPITAL LETTER O WITH BREVE
	335:  "o",     //ŏ	335	014F	 	LATIN SMALL LETTER O WITH BREVE
	336:  "O",     //Ő	336	0150	&Odblac;	LATIN CAPITAL LETTER O WITH DOUBLE ACUTE
	337:  "oe",    //ő	337	0151	&odblac;	LATIN SMALL LETTER O WITH DOUBLE ACUTE
	338:  "Oe",    //Œ	338	0152	&OElig;	LATIN CAPITAL LIGATURE OE
	339:  "oe",    //œ	339	0153	&oelig;	LATIN SMALL LIGATURE OE
	340:  "R",     //Ŕ	340	0154	&Racute;	LATIN CAPITAL LETTER R WITH ACUTE
	341:  "r",     //ŕ	341	0155	&racute;	LATIN SMALL LETTER R WITH ACUTE
	342:  "R",     //Ŗ	342	0156	&Rcedil;	LATIN CAPITAL LETTER R WITH CEDILLA
	343:  "r",     //ŗ	343	0157	&rcedil;	LATIN SMALL LETTER R WITH CEDILLA
	344:  "R",     //Ř	344	0158	&Rcaron;	LATIN CAPITAL LETTER R WITH CARON
	345:  "r",     //ř	345	0159	&rcaron;	LATIN SMALL LETTER R WITH CARON
	346:  "S",     //Ś	346	015A	&Sacute;	LATIN CAPITAL LETTER S WITH ACUTE
	347:  "s",     //ś	347	015B	&sacute;	LATIN SMALL LETTER S WITH ACUTE
	348:  "S",     //Ŝ	348	015C	&Scirc;	LATIN CAPITAL LETTER S WITH CIRCUMFLEX
	349:  "s",     //ŝ	349	015D	&scirc;	LATIN SMALL LETTER S WITH CIRCUMFLEX
	350:  "S",     //Ş	350	015E	&Scedil;	LATIN CAPITAL LETTER S WITH CEDILLA
	351:  "s",     //ş	351	015F	&scedil;	LATIN SMALL LETTER S WITH CEDILLA
	352:  "S",     //Š	352	0160	&Scaron;	LATIN CAPITAL LETTER S WITH CARON
	353:  "s",     //š	353	0161	&scaron;	LATIN SMALL LETTER S WITH CARON
	354:  "T",     //Ţ	354	0162	&Tcedil;	LATIN CAPITAL LETTER T WITH CEDILLA
	355:  "t",     //ţ	355	0163	&tcedil;	LATIN SMALL LETTER T WITH CEDILLA
	356:  "T",     //Ť	356	0164	&Tcaron;	LATIN CAPITAL LETTER T WITH CARON
	357:  "t",     //ť	357	0165	&tcaron;	LATIN SMALL LETTER T WITH CARON
	358:  "T",     //Ŧ	358	0166	&Tstrok;	LATIN CAPITAL LETTER T WITH STROKE
	359:  "t",     //ŧ	359	0167	&tstrok;	LATIN SMALL LETTER T WITH STROKE
	360:  "U",     //Ũ	360	0168	&Utilde;	LATIN CAPITAL LETTER U WITH TILDE
	361:  "u",     //ũ	361	0169	&utilde;	LATIN SMALL LETTER U WITH TILDE
	362:  "U",     //Ū	362	016A	&Umacr;	LATIN CAPITAL LETTER U WITH MACRON
	363:  "u",     //ū	363	016B	&umacr;	LATIN SMALL LETTER U WITH MACRON
	364:  "U",     //Ŭ	364	016C	&Ubreve;	LATIN CAPITAL LETTER U WITH BREVE
	365:  "u",     //ŭ	365	016D	&ubreve;	LATIN SMALL LETTER U WITH BREVE
	366:  "U",     //Ů	366	016E	&Uring;	LATIN CAPITAL LETTER U WITH RING ABOVE
	367:  "u",     //ů	367	016F	&uring;	LATIN SMALL LETTER U WITH RING ABOVE
	368:  "U",     //Ű	368	0170	&Udblac;	LATIN CAPITAL LETTER U WITH DOUBLE ACUTE
	369:  "U",     //ű	369	0171	&udblac;	LATIN SMALL LETTER U WITH DOUBLE ACUTE
	370:  "U",     //Ų	370	0172	&Uogon;	LATIN CAPITAL LETTER U WITH OGONEK
	371:  "u",     //ų	371	0173	&uogon;	LATIN SMALL LETTER U WITH OGONEK
	372:  "W",     //Ŵ	372	0174	&Wcirc;	LATIN CAPITAL LETTER W WITH CIRCUMFLEX
	373:  "w",     //ŵ	373	0175	&wcirc;	LATIN SMALL LETTER W WITH CIRCUMFLEX
	374:  "Y",     //Ŷ	374	0176	&Ycirc;	LATIN CAPITAL LETTER Y WITH CIRCUMFLEX
	375:  "y",     //ŷ	375	0177	&ycirc;	LATIN SMALL LETTER Y WITH CIRCUMFLEX
	376:  "Y",     //Ÿ	376	0178	&Yuml;	LATIN CAPITAL LETTER Y WITH DIAERESIS
	377:  "Z",     //Ź	377	0179	&Zacute;	LATIN CAPITAL LETTER Z WITH ACUTE
	378:  "z",     //ź	378	017A	&zacute;	LATIN SMALL LETTER Z WITH ACUTE
	379:  "Z",     //Ż	379	017B	&Zdot;	LATIN CAPITAL LETTER Z WITH DOT ABOVE
	380:  "z",     // ż
	381:  "Z",     //Ž	381	017D	&Zcaron;	LATIN CAPITAL LETTER Z WITH CARON
	382:  "z",     //ž	382	017E	&zcaron;	LATIN SMALL LETTER Z WITH CARON
	383:  "s",     //ſ	383	017F	 	LATIN SMALL LETTER LONG S
	384:  "b",     //ƀ	384	0180	 	LATIN SMALL LETTER B WITH STROKE
	385:  "B",     //Ɓ	385	0181	 	LATIN CAPITAL LETTER B WITH HOOK
	386:  "B",     //Ƃ	386	0182	 	LATIN CAPITAL LETTER B WITH TOPBAR
	387:  "B",     //	ƃ	387	0183	 	LATIN SMALL LETTER B WITH TOPBAR
	388:  "b",     //Ƅ	388	0184	 	LATIN CAPITAL LETTER TONE SIX
	389:  "b",     //ƅ	389	0185	 	LATIN SMALL LETTER TONE SIX
	390:  "O",     //Ɔ	390	0186	 	LATIN CAPITAL LETTER OPEN O
	391:  "C",     //Ƈ	391	0187	 	LATIN CAPITAL LETTER C WITH HOOK
	392:  "c",     //ƈ	392	0188	 	LATIN SMALL LETTER C WITH HOOK
	393:  "D",     //Ɖ	393	0189	 	LATIN CAPITAL LETTER AFRICAN D
	394:  "D",     //Ɗ	394	018A	 	LATIN CAPITAL LETTER D WITH HOOK
	395:  "D",     //Ƌ	395	018B	 	LATIN CAPITAL LETTER D WITH TOPBAR
	396:  "D",     //ƌ	396	018C	 	LATIN SMALL LETTER D WITH TOPBAR
	397:  "d",     //ƍ	397	018D	 	LATIN SMALL LETTER TURNED DELTA
	398:  "E",     //Ǝ	398	018E	 	LATIN CAPITAL LETTER REVERSED E
	399:  "sch",   //Ə	399	018F	 	LATIN CAPITAL LETTER SCHWA
	400:  "e",     //Ɛ	400	0190	 	LATIN CAPITAL LETTER OPEN E
	401:  "F",     //Ƒ	401	0191	 	LATIN CAPITAL LETTER F WITH HOOK
	402:  "Ff",    //Ƒƒ	402	0192	&fnof;	LATIN SMALL LETTER F WITH HOOK
	403:  "Fg",    //ƑƓ	403	0193	 	LATIN CAPITAL LETTER G WITH HOOK
	404:  "Fy",    //ƑƔ	404	0194	 	LATIN CAPITAL LETTER GAMMA
	405:  "I",     //Ƒƕ	405	0195	 	LATIN SMALL LETTER HV
	406:  "I",     //ƑƖ	406	0196	 	LATIN CAPITAL LETTER IOTA
	407:  "I",     //ƑƗ	407	0197	 	LATIN CAPITAL LETTER I WITH STROKE
	408:  "K",     //ƑƘ	408	0198	 	LATIN CAPITAL LETTER K WITH HOOK
	409:  "k",     //Ƒƙ	409	0199	 	LATIN SMALL LETTER K WITH HOOK
	410:  "l",     //Ƒƚ	410	019A	 	LATIN SMALL LETTER L WITH BAR
	411:  "l",     //Ƒƛ	411	019B	 	LATIN SMALL LETTER LAMBDA WITH STROKE
	412:  "M",     //ƑƜ	412	019C	 	LATIN CAPITAL LETTER TURNED M
	413:  "N",     //ƑƝ	413	019D	 	LATIN CAPITAL LETTER N WITH LEFT HOOK
	414:  "n",     //Ƒƞ	414	019E	 	LATIN SMALL LETTER N WITH LONG RIGHT LEG
	415:  "O",     //ƑƟ	415	019F	 	LATIN CAPITAL LETTER O WITH MIDDLE TILDE
	416:  "O",     //ƑƠ	416	01A0	 	LATIN CAPITAL LETTER O WITH HORN
	417:  "o",     //Ƒơ	417	01A1	 	LATIN SMALL LETTER O WITH HORN
	418:  "Oi",    //ƑƢ	418	01A2	 	LATIN CAPITAL LETTER OI
	419:  "oi",    //Ƒƣ	419	01A3	 	LATIN SMALL LETTER OI
	420:  "P",     //ƑƤ	420	01A4	 	LATIN CAPITAL LETTER P WITH HOOK
	421:  "p",     //Ƒƥ	421	01A5	 	LATIN SMALL LETTER P WITH HOOK
	422:  "yr",    //ƑƦ	422	01A6	 	LATIN LETTER YR
	423:  "",      //ƑƧ	423	01A7	 	LATIN CAPITAL LETTER TONE TWO
	424:  "",      //Ƒƨ	424	01A8	 	LATIN SMALL LETTER TONE TWO
	425:  "",      //ƑƩ	425	01A9	 	LATIN CAPITAL LETTER ESH
	426:  "",      //Ƒƪ	426	01AA	 	LATIN LETTER REVERSED ESH LOOP
	427:  "t",     //Ƒƫ	427	01AB	 	LATIN SMALL LETTER T WITH PALATAL HOOK
	428:  "T",     //ƑƬ	428	01AC	 	LATIN CAPITAL LETTER T WITH HOOK
	429:  "t",     //Ƒƭ	429	01AD	 	LATIN SMALL LETTER T WITH HOOK
	430:  "T",     //ƑƮ	430	01AE	 	LATIN CAPITAL LETTER T WITH RETROFLEX HOOK
	431:  "u",     //ƑƯ	431	01AF	 	LATIN CAPITAL LETTER U WITH HORN
	432:  "u",     //Ƒư	432	01B0	 	LATIN SMALL LETTER U WITH HORN
	433:  "Y",     //ƑƱ	433	01B1	 	LATIN CAPITAL LETTER UPSILON
	434:  "V",     //ƑƲ	434	01B2	 	LATIN CAPITAL LETTER V WITH HOOK
	435:  "Y",     //ƑƳ	435	01B3	 	LATIN CAPITAL LETTER Y WITH HOOK
	436:  "Y",     //Ƒƴ	436	01B4	 	LATIN SMALL LETTER Y WITH HOOK
	437:  "Z",     //ƑƵ	437	01B5	&imped;	LATIN CAPITAL LETTER Z WITH STROKE
	438:  "z",     //Ƒƶ	438	01B6	 	LATIN SMALL LETTER Z WITH STROKE
	439:  "E",     //ƑƷ	439	01B7	 	LATIN CAPITAL LETTER EZH
	440:  "E",     //ƑƸ	440	01B8	 	LATIN CAPITAL LETTER EZH REVERSED
	441:  "e",     //Ƒƹ	441	01B9	 	LATIN SMALL LETTER EZH REVERSED
	442:  "e",     //Ƒƺ	442	01BA	 	LATIN SMALL LETTER EZH WITH TAIL
	443:  "2",     //Ƒƻ	443	01BB	 	LATIN LETTER TWO WITH STROKE
	444:  "5",     //ƑƼ	444	01BC	 	LATIN CAPITAL LETTER TONE FIVE
	445:  "5",     //Ƒƽ	445	01BD	 	LATIN SMALL LETTER TONE FIVE
	446:  "",      //Ƒƾ	446	01BE	 	LATIN LETTER INVERTED GLOTTAL STOP WITH STROKE
	447:  "",      //Ƒƿ	447	01BF	 	LATIN LETTER WYNN
	448:  "",      //Ƒǀ	448	01C0	 	LATIN LETTER DENTAL CLICK
	449:  "",      //Ƒǁ	449	01C1	 	LATIN LETTER LATERAL CLICK
	450:  "",      //Ƒǂ	450	01C2	 	LATIN LETTER ALVEOLAR CLICK
	451:  "",      //Ƒǃ	451	01C3	 	LATIN LETTER RETROFLEX CLICK
	452:  "Dz",    //ƑǄ	452	01C4	 	LATIN CAPITAL LETTER DZ WITH CARON
	453:  "D",     //Ƒǅ	453	01C5	 	LATIN CAPITAL LETTER D WITH SMALL LETTER Z WITH CARON
	454:  "dz",    //Ƒǆ	454	01C6	 	LATIN SMALL LETTER DZ WITH CARON
	455:  "Lj",    //ƑǇ	455	01C7	 	LATIN CAPITAL LETTER LJ
	456:  "Lj",    //Ƒǈ	456	01C8	 	LATIN CAPITAL LETTER L WITH SMALL LETTER J
	457:  "lj",    //Ƒǉ	457	01C9	 	LATIN SMALL LETTER LJ
	458:  "Nj",    //ƑǊ	458	01CA	 	LATIN CAPITAL LETTER NJ
	459:  "Nj",    //Ƒǋ	459	01CB	 	LATIN CAPITAL LETTER N WITH SMALL LETTER J
	460:  "nj",    //Ƒǌ	460	01CC	 	LATIN SMALL LETTER NJ
	461:  "A",     //ƑǍ	461	01CD	 	LATIN CAPITAL LETTER A WITH CARON
	462:  "a",     //Ƒǎ	462	01CE	 	LATIN SMALL LETTER A WITH CARON
	463:  "I",     //ƑǏ	463	01CF	 	LATIN CAPITAL LETTER I WITH CARON
	464:  "i",     //Ƒǐ	464	01D0	 	LATIN SMALL LETTER I WITH CARON
	465:  "O",     //ƑǑ	465	01D1	 	LATIN CAPITAL LETTER O WITH CARON
	466:  "o",     //Ƒǒ	466	01D2	 	LATIN SMALL LETTER O WITH CARON
	467:  "U",     //ƑǓ	467	01D3	 	LATIN CAPITAL LETTER U WITH CARON
	468:  "u",     //Ƒǔ	468	01D4	 	LATIN SMALL LETTER U WITH CARON
	469:  "U",     //ƑǕ	469	01D5	 	LATIN CAPITAL LETTER U WITH DIAERESIS AND MACRON
	470:  "u",     //Ƒǖ	470	01D6	 	LATIN SMALL LETTER U WITH DIAERESIS AND MACRON
	471:  "U",     //ƑǗ	471	01D7	 	LATIN CAPITAL LETTER U WITH DIAERESIS AND ACUTE
	472:  "u",     //Ƒǘ	472	01D8	 	LATIN SMALL LETTER U WITH DIAERESIS AND ACUTE
	473:  "U",     //ƑǙ	473	01D9	 	LATIN CAPITAL LETTER U WITH DIAERESIS AND CARON
	474:  "u",     //Ƒǚ	474	01DA	 	LATIN SMALL LETTER U WITH DIAERESIS AND CARON
	475:  "U",     //ƑǛ	475	01DB	 	LATIN CAPITAL LETTER U WITH DIAERESIS AND GRAVE
	476:  "u",     //Ƒǜ	476	01DC	 	LATIN SMALL LETTER U WITH DIAERESIS AND GRAVE
	477:  "e",     //Ƒǝ	477	01DD	 	LATIN SMALL LETTER TURNED E
	478:  "A",     //ƑǞ	478	01DE	 	LATIN CAPITAL LETTER A WITH DIAERESIS AND MACRON
	479:  "a",     //Ƒǟ	479	01DF	 	LATIN SMALL LETTER A WITH DIAERESIS AND MACRON
	480:  "A",     //ƑǠ	480	01E0	 	LATIN CAPITAL LETTER A WITH DOT ABOVE AND MACRON
	481:  "a",     //Ƒǡ	481	01E1	 	LATIN SMALL LETTER A WITH DOT ABOVE AND MACRON
	482:  "Ae",    //ƑǢ	482	01E2	 	LATIN CAPITAL LETTER AE WITH MACRON
	483:  "ae",    //Ƒǣ	483	01E3	 	LATIN SMALL LETTER AE WITH MACRON
	484:  "G",     //ƑǤ	484	01E4	 	LATIN CAPITAL LETTER G WITH STROKE
	485:  "g",     //Ƒǥ	485	01E5	 	LATIN SMALL LETTER G WITH STROKE
	486:  "G",     //ƑǦ	486	01E6	 	LATIN CAPITAL LETTER G WITH CARON
	487:  "g",     //Ƒǧ	487	01E7	 	LATIN SMALL LETTER G WITH CARON
	488:  "K",     //ƑǨ	488	01E8	 	LATIN CAPITAL LETTER K WITH CARON
	489:  "k",     //Ƒǩ	489	01E9	 	LATIN SMALL LETTER K WITH CARON
	490:  "O",     // Ǫ	490	01EA	 	LATIN CAPITAL LETTER O WITH OGONEK
	491:  "o",     //ǫ	491	01EB	 	LATIN SMALL LETTER O WITH OGONEK
	492:  "O",     //Ǭ	492	01EC	 	LATIN CAPITAL LETTER O WITH OGONEK AND MACRON
	493:  "o",     //ǭ	493	01ED	 	LATIN SMALL LETTER O WITH OGONEK AND MACRON
	494:  "Ezh",   //Ǯ	494	01EE	 	LATIN CAPITAL LETTER EZH WITH CARON
	495:  "esz",   //ǯ	495	01EF	 	LATIN SMALL LETTER EZH WITH CARON
	496:  "j",     //ǰ	496	01F0	 	LATIN SMALL LETTER J WITH CARON
	497:  "Dz",    //Ǳ	497	01F1	 	LATIN CAPITAL LETTER DZ
	498:  "Dz",    //ǲ	498	01F2	 	LATIN CAPITAL LETTER D WITH SMALL LETTER Z
	499:  "dz",    //ǳ	499	01F3	 	LATIN SMALL LETTER DZ
	500:  "G",     //Ǵ	500	01F4	 	LATIN CAPITAL LETTER G WITH ACUTE
	501:  "g",     //ǵ	501	01F5	&gacute;	LATIN SMALL LETTER G WITH ACUTE
	502:  "Hwair", //Ƕ	502	01F6	 	LATIN CAPITAL LETTER HWAIR
	503:  "Wynn",  //Ƿ	503	01F7	 	LATIN CAPITAL LETTER WYNN
	504:  "N",     //Ǹ	504	01F8	 	LATIN CAPITAL LETTER N WITH GRAVE
	505:  "n",     //ǹ	505	01F9	 	LATIN SMALL LETTER N WITH GRAVE
	506:  "A",     //Ǻ	506	01FA	 	LATIN CAPITAL LETTER A WITH RING ABOVE AND ACUTE
	507:  "a",     //ǻ	507	01FB	 	LATIN SMALL LETTER A WITH RING ABOVE AND ACUTE
	508:  "Ae",    //Ǽ	508	01FC	 	LATIN CAPITAL LETTER AE WITH ACUTE
	509:  "ae",    //ǽ	509	01FD	 	LATIN SMALL LETTER AE WITH ACUTE
	510:  "O",     //Ǿ	510	01FE	 	LATIN CAPITAL LETTER O WITH STROKE AND ACUTE
	511:  "o",     //ǿ	511	01FF	 	LATIN SMALL LETTER O WITH STROKE AND ACUTE
	512:  "A",     //Ȁ	512	0200	 	LATIN CAPITAL LETTER A WITH DOUBLE GRAVE
	513:  "a",     //ȁ	513	0201	 	LATIN SMALL LETTER A WITH DOUBLE GRAVE
	514:  "A",     //Ȃ	514	0202	 	LATIN CAPITAL LETTER A WITH INVERTED BREVE
	515:  "a",     //ȃ	515	0203	 	LATIN SMALL LETTER A WITH INVERTED BREVE
	516:  "E",     //Ȅ	516	0204	 	LATIN CAPITAL LETTER E WITH DOUBLE GRAVE
	517:  "e",     //ȅ	517	0205	 	LATIN SMALL LETTER E WITH DOUBLE GRAVE
	518:  "E",     //Ȇ	518	0206	 	LATIN CAPITAL LETTER E WITH INVERTED BREVE
	519:  "e",     //ȇ	519	0207	 	LATIN SMALL LETTER E WITH INVERTED BREVE
	520:  "E",     //Ȉ	520	0208	 	LATIN CAPITAL LETTER I WITH DOUBLE GRAVE
	521:  "i",     //ȉ	521	0209	 	LATIN SMALL LETTER I WITH DOUBLE GRAVE
	522:  "I",     //Ȋ	522	020A	 	LATIN CAPITAL LETTER I WITH INVERTED BREVE
	523:  "i",     //ȋ	523	020B	 	LATIN SMALL LETTER I WITH INVERTED BREVE
	524:  "O",     //Ȍ	524	020C	 	LATIN CAPITAL LETTER O WITH DOUBLE GRAVE
	525:  "o",     //ȍ	525	020D	 	LATIN SMALL LETTER O WITH DOUBLE GRAVE
	526:  "O",     //Ȏ	526	020E	 	LATIN CAPITAL LETTER O WITH INVERTED BREVE
	527:  "o",     //ȏ	527	020F	 	LATIN SMALL LETTER O WITH INVERTED BREVE
	528:  "R",     //Ȑ	528	0210	 	LATIN CAPITAL LETTER R WITH DOUBLE GRAVE
	529:  "r",     //ȑ	529	0211	 	LATIN SMALL LETTER R WITH DOUBLE GRAVE
	530:  "R",     //Ȓ	530	0212	 	LATIN CAPITAL LETTER R WITH INVERTED BREVE
	531:  "r",     //ȓ	531	0213	 	LATIN SMALL LETTER R WITH INVERTED BREVE
	532:  "U",     //Ȕ	532	0214	 	LATIN CAPITAL LETTER U WITH DOUBLE GRAVE
	533:  "u",     //ȕ	533	0215	 	LATIN SMALL LETTER U WITH DOUBLE GRAVE
	534:  "U",     //Ȗ	534	0216	 	LATIN CAPITAL LETTER U WITH INVERTED BREVE
	535:  "u",     //ȗ	535	0217	 	LATIN SMALL LETTER U WITH INVERTED BREVE
	536:  "S",     //Ș	536	0218	 	LATIN CAPITAL LETTER S WITH COMMA BELOW
	537:  "s",     //ș	537	0219	 	LATIN SMALL LETTER S WITH COMMA BELOW
	538:  "T",     //Ț	538	021A	 	LATIN CAPITAL LETTER T WITH COMMA BELOW
	539:  "t",     //ț	539	021B	 	LATIN SMALL LETTER T WITH COMMA BELOW
	540:  "z",     //Ȝ	540	021C	 	LATIN CAPITAL LETTER YOGH
	541:  "z",     //ȝ	541	021D	 	LATIN SMALL LETTER YOGH
	542:  "H",     //Ȟ	542	021E	 	LATIN CAPITAL LETTER H WITH CARON
	543:  "h",     //ȟ	543	021F	 	LATIN SMALL LETTER H WITH CARON
	544:  "n",     //Ƞ	544	0220	 	LATIN CAPITAL LETTER N WITH LONG RIGHT LEG
	545:  "d",     //ȡ	545	0221	 	LATIN SMALL LETTER D WITH CURL
	546:  "Ou",    //Ȣ	546	0222	 	LATIN CAPITAL LETTER OU
	547:  "ou",    //ȣ	547	0223	 	LATIN SMALL LETTER OU
	548:  "Z",     //Ȥ	548	0224	 	LATIN CAPITAL LETTER Z WITH HOOK
	549:  "z",     //ȥ	549	0225	 	LATIN SMALL LETTER Z WITH HOOK
	550:  "A",     //Ȧ	550	0226	 	LATIN CAPITAL LETTER A WITH DOT ABOVE
	551:  "a",     //ȧ	551	0227	 	LATIN SMALL LETTER A WITH DOT ABOVE
	552:  "E",     //Ȩ	552	0228	 	LATIN CAPITAL LETTER E WITH CEDILLA
	553:  "e",     //ȩ	553	0229	 	LATIN SMALL LETTER E WITH CEDILLA
	554:  "O",     //Ȫ	554	022A	 	LATIN CAPITAL LETTER O WITH DIAERESIS AND MACRON
	555:  "o",     //ȫ	555	022B	 	LATIN SMALL LETTER O WITH DIAERESIS AND MACRON
	556:  "O",     //Ȭ	556	022C	 	LATIN CAPITAL LETTER O WITH TILDE AND MACRON
	557:  "o",     //ȭ	557	022D	 	LATIN SMALL LETTER O WITH TILDE AND MACRON
	558:  "O",     //Ȯ	558	022E	 	LATIN CAPITAL LETTER O WITH DOT ABOVE
	559:  "o",     //ȯ	559	022F	 	LATIN SMALL LETTER O WITH DOT ABOVE
	560:  "O",     //Ȱ	560	0230	 	LATIN CAPITAL LETTER O WITH DOT ABOVE AND MACRON
	561:  "o",     //ȱ	561	0231	 	LATIN SMALL LETTER O WITH DOT ABOVE AND MACRON
	562:  "Y",     //Ȳ	562	0232	 	LATIN CAPITAL LETTER Y WITH MACRON
	563:  "y",     //ȳ	563	0233	 	LATIN SMALL LETTER Y WITH MACRON
	564:  "l",     //ȴ	564	0234	 	LATIN SMALL LETTER L WITH CURL
	565:  "n",     //ȵ	565	0235	 	LATIN SMALL LETTER N WITH CURL
	566:  "t",     //ȶ	566	0236	 	LATIN SMALL LETTER T WITH CURL
	567:  "j",     //ȷ	567	0237	&jmath;	LATIN SMALL LETTER DOTLESS J
	568:  "db",    //ȸ	568	0238	 	LATIN SMALL LETTER DB DIGRAPH
	569:  "qp",    //ȹ	569	0239	 	LATIN SMALL LETTER QP DIGRAPH
	570:  "A",     //Ⱥ	570	023A	 	LATIN CAPITAL LETTER A WITH STROKE
	571:  "C",     //Ȼ	571	023B	 	LATIN CAPITAL LETTER C WITH STROKE
	572:  "c",     //ȼ	572	023C	 	LATIN SMALL LETTER C WITH STROKE
	573:  "L",     //Ƚ	573	023D	 	LATIN CAPITAL LETTER L WITH BAR
	574:  "T",     //Ⱦ	574	023E	 	LATIN CAPITAL LETTER T WITH DIAGONAL STROKE
	575:  "s",     //ȿ	575	023F	 	LATIN SMALL LETTER S WITH SWASH TAIL
	576:  "z",     //ɀ	576	0240	 	LATIN SMALL LETTER Z WITH SWASH TAIL
	579:  "B",     //Ƀ	579	0243	 	LATIN CAPITAL LETTER B WITH STROKE
	580:  "U",     //Ʉ	580	0244	 	LATIN CAPITAL LETTER U BAR
	582:  "E",     //Ɇ	582	0246	 	LATIN CAPITAL LETTER E WITH STROKE
	583:  "e",     //ɇ	583	0247	 	LATIN SMALL LETTER E WITH STROKE
	584:  "J",     //Ɉ	584	0248	 	LATIN CAPITAL LETTER J WITH STROKE
	585:  "j",     //ɉ	585	0249	 	LATIN SMALL LETTER J WITH STROKE
	586:  "q",     //Ɋ	586	024A	 	LATIN CAPITAL LETTER SMALL Q WITH HOOK TAIL
	587:  "q",     //ɋ	587	024B	 	LATIN SMALL LETTER Q WITH HOOK TAIL
	588:  "R",     //Ɍ	588	024C	 	LATIN CAPITAL LETTER R WITH STROKE
	589:  "r",     //ɍ	589	024D	 	LATIN SMALL LETTER R WITH STROKE
	590:  "Y",     //Ɏ	590	024E	 	LATIN CAPITAL LETTER Y WITH STROKE
	591:  "y",     //ɏ	591	024F
	688:  "h",     //ʰ	688	02B0	 	MODIFIER LETTER SMALL H
	689:  "h",     //ʱ	689	02B1	 	MODIFIER LETTER SMALL H WITH HOOK
	690:  "j",     //ʲ	690	02B2	 	MODIFIER LETTER SMALL J
	691:  "r",     //ʳ	691	02B3	 	MODIFIER LETTER SMALL R
	692:  "r",     //ʴ	692	02B4	 	MODIFIER LETTER SMALL TURNED R
	693:  "r",     //ʵ	693	02B5	 	MODIFIER LETTER SMALL TURNED R WITH HOOK
	694:  "R",     //ʶ	694	02B6	 	MODIFIER LETTER SMALL CAPITAL INVERTED R
	695:  "w",     //ʷ	695	02B7	 	MODIFIER LETTER SMALL W
	696:  "y",     //ʸ	696	02B8	 	MODIFIER LETTER SMALL Y
	737:  "l",     //ˡ		737	02E1	 	MODIFIER LETTER SMALL L
	738:  "s",     //	ˢ	738	02E2	 	MODIFIER LETTER SMALL S
	739:  "x",     //ˣ		739	02E3	 	MODIFIER LETTER SMALL X
	768:  "o",     //ò	768	0300	 	GRAVE ACCENT
	769:  "o",     //ó	769	0301	 	ACUTE ACCENT
	770:  "o",     //ô	770	0302	 	CIRCUMFLEX ACCENT
	771:  "",      //õ	771	0303	 	TILDE
	772:  "o",     //ō	772	0304	 	MACRON
	773:  "o",     //o̅	773	0305	 	OVERLINE
	774:  "o",     //ŏ	774	0306	 	BREVE
	775:  "o",     //ȯ	775	0307	 	DOT ABOVE
	776:  "o",     //ö	776	0308	 	DIAERESIS
	777:  "o",     //ỏ	777	0309	 	HOOK ABOVE
	778:  "o",     //o̊	778	030A	 	RING ABOVE
	779:  "o",     //ő	779	030B	 	DOUBLE ACUTE ACCENT
	780:  "o",     //ǒ	780	030C	 	CARON
	781:  "o",     //o̍	781	030D	 	VERTICAL LINE ABOVE
	782:  "o",     //o̎	782	030E	 	DOUBLE VERTICAL LINE ABOVE
	783:  "o",     //ȍ	783	030F	 	DOUBLE GRAVE ACCENT
	784:  "o",     //o̐	784	0310	 	CANDRABINDU
	785:  "o",     //ȏ	785	0311	 	INVERTED BREVE
	786:  "o",     //o̒	786	0312	 	TURNED COMMA ABOVE
	787:  "o",     //o̓	787	0313	 	COMMA ABOVE
	788:  "o",     //o̔	788	0314	 	REVERSED COMMA ABOVE
	789:  "o",     //o̕	789	0315	 	COMMA ABOVE RIGHT
	790:  "o",     //o̖	790	0316	 	GRAVE ACCENT BELOW
	791:  "o",     //o̗	791	0317	 	ACUTE ACCENT BELOW
	792:  "o",     //o̘	792	0318	 	LEFT TACK BELOW
	793:  "o",     //o̙	793	0319	 	RIGHT TACK BELOW
	794:  "o",     //o̚	794	031A	 	LEFT ANGLE ABOVE
	795:  "o",     //ơ	795	031B	 	HORN
	796:  "o",     //o̜	796	031C	 	LEFT HALF RING BELOW
	797:  "o",     //o̝	797	031D	 	UP TACK BELOW
	798:  "o",     //o̞	798	031E	 	DOWN TACK BELOW
	799:  "o",     //o̟	799	031F	 	PLUS SIGN BELOW
	800:  "o",     //o̠	800	0320	 	MINUS SIGN BELOW
	801:  "o",     //o̡	801	0321	 	PALATALIZED HOOK BELOW
	802:  "o",     //o̢	802	0322	 	RETROFLEX HOOK BELOW
	803:  "o",     //ọ	803	0323	 	DOT BELOW
	804:  "o",     //o̤	804	0324	 	DIAERESIS BELOW
	805:  "o",     //o̥	805	0325	 	RING BELOW
	806:  "o",     //o̦	806	0326	 	COMMA BELOW
	807:  "o",     //o̧	807	0327	 	CEDILLA
	808:  "o",     //ǫ	808	0328	 	OGONEK
	809:  "o",     //o̩	809	0329	 	VERTICAL LINE BELOW
	810:  "o",     //o̪	810	032A	 	BRIDGE BELOW
	811:  "o",     //o̫	811	032B	 	INVERTED DOUBLE ARCH BELOW
	812:  "o",     //o̬	812	032C	 	CARON BELOW
	813:  "o",     //o̭	813	032D	 	CIRCUMFLEX ACCENT BELOW
	814:  "o",     //o̮	814	032E	 	BREVE BELOW
	815:  "o",     //o̯	815	032F	 	INVERTED BREVE BELOW
	816:  "o",     //o̰	816	0330	 	TILDE BELOW
	817:  "o",     //o̱	817	0331	 	MACRON BELOW
	818:  "o",     //o̲	818	0332	 	LOW LINE
	819:  "o",     //o̳	819	0333	 	DOUBLE LOW LINE
	820:  "o",     //o̴	820	0334	 	TILDE OVERLAY
	821:  "o",     //o̵	821	0335	 	SHORT STROKE OVERLAY
	822:  "o",     //o̶	822	0336	 	LONG STROKE OVERLAY
	823:  "o",     //o̷	823	0337	 	SHORT SOLIDUS OVERLAY
	824:  "o",     //o̸	824	0338	 	LONG SOLIDUS OVERLAY
	825:  "o",     //o̹	825	0339	 	RIGHT HALF RING BELOW
	826:  "o",     //o̺	826	033A	 	INVERTED BRIDGE BELOW
	827:  "o",     //o̻	827	033B	 	SQUARE BELOW
	828:  "o",     //o̼	828	033C	 	SEAGULL BELOW
	829:  "o",     //o̽	829	033D	 	X ABOVE
	830:  "o",     //o̾	830	033E	 	VERTICAL TILDE
	831:  "o",     //o̿	831	033F	 	DOUBLE OVERLINE
	832:  "o",     //ò	832	0340	 	GRAVE TONE MARK
	833:  "o",     //ó	833	0341	 	ACUTE TONE MARK
	834:  "o",     //o͂	834	0342	 	GREEK PERISPOMENI (combined with theta)
	835:  "o",     //o̓	835	0343	 	GREEK KORONIS (combined with theta)
	836:  "o",     //ö́	836	0344	 	GREEK DIALYTIKA TONOS (combined with theta)
	837:  "o",     //oͅ	837	0345	 	GREEK YPOGEGRAMMENI (combined with theta)
	838:  "o",     //o͆	838	0346	 	BRIDGE ABOVE
	839:  "o",     //o͇	839	0347	 	EQUALS SIGN BELOW
	840:  "o",     //o͈	840	0348	 	DOUBLE VERTICAL LINE BELOW
	841:  "o",     //o͉	841	0349	 	LEFT ANGLE BELOW
	842:  "o",     //	o͊	842	034A	 	NOT TILDE ABOVE
	843:  "o",     //	o͋	843	034B	 	HOMOTHETIC ABOVE
	844:  "o",     //o͌	844	034C	 	ALMOST EQUAL TO ABOVE
	845:  "o",     //o͍	845	034D	 	LEFT RIGHT ARROW BELOW
	846:  "o",     //o͎	846	034E	 	UPWARDS ARROW BELOW
	847:  "o",     //o͏	847	034F	 	GRAPHEME JOINER
	848:  "o",     //o͐	848	0350	 	RIGHT ARROWHEAD ABOVE
	849:  "o",     //o͑	849	0351	 	LEFT HALF RING ABOVE
	850:  "o",     //o͒	850	0352	 	FERMATA
	851:  "o",     //o͓	851	0353	 	X BELOW
	852:  "o",     //o͔	852	0354	 	LEFT ARROWHEAD BELOW
	853:  "o",     //o͕	853	0355	 	RIGHT ARROWHEAD BELOW
	854:  "o",     //o͖	854	0356	 	RIGHT ARROWHEAD AND UP ARROWHEAD BELOW
	855:  "o",     //o͗	855	0357	 	RIGHT HALF RING ABOVE
	856:  "o",     //o͘	856	0358	 	DOT ABOVE RIGHT
	857:  "o",     //o͙	857	0359	 	ASTERISK BELOW
	858:  "o",     //o͚	858	035A	 	DOUBLE RING BELOW
	859:  "o",     //	o͛	859	035B	 	ZIGZAG ABOVE
	860:  "o",     //	͜o	860	035C	 	DOUBLE BREVE BELOW
	861:  "o",     //	͝o	861	035D	 	DOUBLE BREVE
	862:  "o",     //	͞o	862	035E	 	DOUBLE MACRON
	863:  "o",     //	͟o	863	035F	 	DOUBLE MACRON BELOW
	864:  "o",     //	͠o	864	0360	 	DOUBLE TILDE
	865:  "o",     //	͡o	865	0361	 	DOUBLE INVERTED BREVE
	866:  "o",     //	͢o	866	0362	 	DOUBLE RIGHTWARDS ARROW BELOW
	867:  "o",     //	oͣ	867	0363	 	LATIN SMALL LETTER A
	868:  "e",     //oͤ	868	0364	 	LATIN SMALL LETTER E
	869:  "i",     //oͥ	869	0365	 	LATIN SMALL LETTER I
	870:  "o",     //oͦ	870	0366	 	LATIN SMALL LETTER O
	871:  "u",     //oͧ	871	0367	 	LATIN SMALL LETTER U
	872:  "c",     //oͨ	872	0368	 	LATIN SMALL LETTER C
	873:  "d",     //oͩ	873	0369	 	LATIN SMALL LETTER D
	874:  "h",     //oͪ	874	036A	 	LATIN SMALL LETTER H
	875:  "m",     //oͫ	875	036B	 	LATIN SMALL LETTER M
	876:  "r",     //oͬ	876	036C	 	LATIN SMALL LETTER R
	877:  "t",     //oͭ	877	036D	 	LATIN SMALL LETTER T
	878:  "v",     //oͮ	878	036E	 	LATIN SMALL LETTER V
	879:  "x",     //oͯ	879	036F	 	LATIN SMALL LETTER X
	1104: "e",     //	ѐ	1104	0450	 	CYRILLIC SMALL LETTER IE WITH GRAVE
	1105: "e",     //	ё	1105	0451	 	CYRILLIC SMALL LETTER IO
	1106: "dje",   //	ђ	1106	0452	 	CYRILLIC SMALL LETTER DJE
	1107: "gje",   //	ѓ	1107	0453	 	CYRILLIC SMALL LETTER GJE
	1108: "ie",    //	є	1108	0454	 	CYRILLIC SMALL LETTER UKRAINIAN IE
	1109: "dze",   //	ѕ	1109	0455	 	CYRILLIC SMALL LETTER DZE
	1110: "i",     //	і	1110	0456	 	CYRILLIC SMALL LETTER BYELORUSSIAN-UKRAINIAN I
	1111: "yi",    //	ї	1111	0457	 	CYRILLIC SMALL LETTER YI
	1112: "je",    //	ј	1112	0458	 	CYRILLIC SMALL LETTER JE
	1113: "lje",   //	љ	1113	0459	 	CYRILLIC SMALL LETTER LJE
	1114: "nje",   //	њ	1114	045A	 	CYRILLIC SMALL LETTER NJE
	1115: "tshe",  //	ћ	1115	045B	 	CYRILLIC SMALL LETTER TSHE
	1116: "kje",   //	ќ	1116	045C	 	CYRILLIC SMALL LETTER KJE
	1117: "i",     //	ѝ	1117	045D	 	CYRILLIC SMALL LETTER I WITH GRAVE
	7840: "A",     // Ạ
	7841: "a",     //ạ
	7842: "A",     // Ả
	7843: "a",     //ả
	7844: "A",     // Ấ
	7845: "a",     // ấ 7845
	7846: "A",     // Ầ
	7847: "a",     //  ầ
	7848: "A",     // Ẩ
	7849: "a",     //ẩ
	7850: "A",     //Ẫ
	7879: "e",     // ệ 7879
	7891: "o",     // ồ
	8357: "m",     // ₥
}

var CharactersToChange = map[int]string{
	0:    "0",     //NUL	0	0000	null character
	1:    "1",     //SOH	1	0001	start of header
	2:    "2",     //STX	2	0002	start of text
	3:    "3",     //ETX	3	0003	end of text
	4:    "4",     //EOT	4	0004	end of transmission
	5:    "",      //ENQ	5	0005	enquiry
	6:    "",      //ACK	6	0006	acknowledge
	7:    "",      //BEL	7	0007	bell (ring)
	8:    "",      //BS	8	0008	backspace
	9:    "",      //HT	9	0009	horizontal tab
	10:   "",      //LF	10	000A	line feed
	11:   "",      //VT	11	000B	vertical tab
	12:   "",      //FF	12	000C	form feed
	13:   "",      //CR	13	000D	carriage return
	14:   "",      //SO	14	000E	shift out
	15:   "",      //SI	15	000F	shift in
	16:   "",      //DLE	16	0010	data link escape
	17:   "",      //DC1	17	0011	device control 1
	18:   "",      //DC2	18	0012	device control 2
	19:   "",      //DC3	19	0013	device control 3
	20:   "",      //DC4	20	0014	device control 4
	21:   "",      //NAK	21	0015	negative acknowledge
	22:   "",      //SYN	22	0016	synchronize
	23:   "",      //ETB	23	0017	end transmission block
	24:   "",      //CAN	24	0018	cancel
	25:   "",      //EM	25	0019	end of medium
	26:   "",      //SUB	26	001A	substitute
	27:   "",      //ESC	27	001B	escape
	28:   "",      //FS	28	001C	file separator
	29:   "",      //GS	29	001D	group separator
	30:   "",      //RS	30	001E	record separator
	31:   "",      //US	31	001F	unit separator
	32:   " ",     //	32	0020	 	SPACE
	33:   "",      //!	33	0021	 	EXCLAMATION MARK
	34:   "",      //"	34	0022	&quot;	QUOTATION MARK
	35:   "",      //#	35	0023	 	NUMBER SIGN
	36:   "",      //$	36	0024	 	DOLLAR SIGN
	37:   "",      //%	37	0025	 	PERCENT SIGN
	38:   "",      //&	38	0026	&amp;	AMPERSAND
	39:   "",      //'	39	0027	 	APOSTROPHE
	40:   "",      //(	40	0028	 	LEFT PARENTHESIS
	41:   "",      //)	41	0029	 	RIGHT PARENTHESIS
	42:   "",      //*	42	002A	 	ASTERISK
	43:   "",      //+	43	002B	 	PLUS SIGN
	44:   "",      //,	44	002C	 	COMMA
	45:   "",      //-	45	002D	 	HYPHEN-MINUS
	46:   ".",     //.	46	002E	 	FULL STOP
	47:   "",      ///	47	002F	 	SOLIDUS
	48:   "0",     //0	48	0030	 	DIGIT ZERO
	49:   "1",     //1	49	0031	 	DIGIT ONE
	50:   "2",     //2	50	0032	 	DIGIT TWO
	51:   "3",     //3 0033	 	DIGIT THREE
	52:   "4",     //4 0034	 	DIGIT FOUR
	53:   "5",     //5 0035	 	DIGIT FIVE
	54:   "6",     //6 0036	 	DIGIT SIX
	55:   "7",     //7 0037	 	DIGIT SEVEN
	56:   "8",     //8 0038	 	DIGIT EIGHT
	57:   "9",     //9 0039	 	DIGIT NINE
	58:   "",      //:	003A	 	COLON
	59:   "",      //;	003B	 	SEMICOLON
	60:   "",      //<	003C	&lt;	LESS-THAN SIGN
	61:   "",      //=	003D	 	EQUALS SIGN
	62:   "",      //>	003E	&gt;	GREATER-THAN SIGN
	63:   "",      //?	003F	 	QUESTION MARK
	64:   "",      //@	0040	 	COMMERCIAL AT
	65:   "A",     //A	0041	 	LATIN CAPITAL LETTER A
	66:   "B",     //B	0042	 	LATIN CAPITAL LETTER B
	67:   "C",     //C	0043	 	LATIN CAPITAL LETTER C
	68:   "D",     //D	0044	 	LATIN CAPITAL LETTER D
	69:   "E",     //E 0045	 	LATIN CAPITAL LETTER E
	70:   "F",     //F	70	0046	 	LATIN CAPITAL LETTER F
	71:   "G",     //G	71	0047	 	LATIN CAPITAL LETTER G
	72:   "H",     //H	72	0048	 	LATIN CAPITAL LETTER H
	73:   "I",     //I	73	0049	 	LATIN CAPITAL LETTER I
	74:   "J",     //J	74	004A	 	LATIN CAPITAL LETTER J
	75:   "K",     //K	75	004B	 	LATIN CAPITAL LETTER K
	76:   "L",     //L 76	004C	 	LATIN CAPITAL LETTER L
	77:   "M",     //M	77	004D	 	LATIN CAPITAL LETTER M
	78:   "N",     //N	78	004E	 	LATIN CAPITAL LETTER N
	79:   "O",     //O	79	004F	 	LATIN CAPITAL LETTER O
	80:   "P",     //P	80	0050	 	LATIN CAPITAL LETTER P
	81:   "Q",     //Q	81	0051	 	LATIN CAPITAL LETTER Q
	82:   "R",     //R	82	0052	 	LATIN CAPITAL LETTER R
	83:   "S",     //S	83	0053	 	LATIN CAPITAL LETTER S
	84:   "T",     //T	84	0054	 	LATIN CAPITAL LETTER T
	85:   "U",     //U	85	0055	 	LATIN CAPITAL LETTER U
	86:   "V",     //V	86	0056	 	LATIN CAPITAL LETTER V
	87:   "W",     //W	87	0057	 	LATIN CAPITAL LETTER W
	88:   "X",     //X	88	0058	 	LATIN CAPITAL LETTER X
	89:   "Y",     //Y	89	0059	 	LATIN CAPITAL LETTER Y
	90:   "Z",     //Z	90	005A	 	LATIN CAPITAL LETTER Z
	91:   "",      // [	91	005B	 	LEFT SQUARE BRACKET
	92:   "",      // \	92	005C	 	REVERSE SOLIDUS
	93:   "",      //]	93	005D	 	RIGHT SQUARE BRACKET
	94:   "",      // ^	94	005E	 	CIRCUMFLEX ACCENT
	95:   "",      //_	95	005F	 	LOW LINE
	96:   "",      //`	96	0060	 	GRAVE ACCENT
	97:   "a",     // a	97	0061	 	LATIN SMALL LETTER A
	98:   "b",     // b	98	0062	 	LATIN SMALL LETTER B
	99:   "c",     //c	99	0063	 	LATIN SMALL LETTER C
	100:  "d",     //d	100	0064	 	LATIN SMALL LETTER D
	101:  "e",     //e	101	0065	 	LATIN SMALL LETTER E
	102:  "f",     //f	102	0066	 	LATIN SMALL LETTER F
	103:  "g",     //g	103	0067	 	LATIN SMALL LETTER G
	104:  "h",     //h	104	0068	 	LATIN SMALL LETTER H
	105:  "i",     //i	105	0069	 	LATIN SMALL LETTER I
	106:  "j",     //j	106	006A	 	LATIN SMALL LETTER J
	107:  "k",     //k	107	006B	 	LATIN SMALL LETTER K
	108:  "l",     //l	108	006C	 	LATIN SMALL LETTER L
	109:  "m",     //m	109	006D	 	LATIN SMALL LETTER M
	110:  "n",     //n	110	006E	 	LATIN SMALL LETTER N
	111:  "o",     //o	111	006F	 	LATIN SMALL LETTER O
	112:  "p",     //p	112	0070	 	LATIN SMALL LETTER P
	113:  "q",     //q	113	0071	 	LATIN SMALL LETTER Q
	114:  "r",     //r	114	0072	 	LATIN SMALL LETTER R
	115:  "s",     //s	115	0073	 	LATIN SMALL LETTER S
	116:  "t",     //t	116	0074	 	LATIN SMALL LETTER T
	117:  "u",     // u	117	0075	 	LATIN SMALL LETTER U
	118:  "v",     //v	118	0076	 	LATIN SMALL LETTER V
	119:  "w",     //w	119	0077	 	LATIN SMALL LETTER W
	120:  "x",     //x	120	0078	 	LATIN SMALL LETTER X
	121:  "y",     //y	121	0079	 	LATIN SMALL LETTER Y
	122:  "z",     //z	122	007A	 	LATIN SMALL LETTER Z
	123:  "",      //{	123	007B	 	LEFT CURLY BRACKET
	124:  "",      //|	124	007C	 	VERTICAL LINE
	125:  "",      //}	125	007D	 	RIGHT CURLY BRACKET
	126:  "",      //~	126	007E	 	TILDE
	127:  "",      // DEL	127	007F	delete (rubout)
	128:  "",      //€	128	0080	CONTROL
	129:  "",      // 	129	0081	CONTROL
	130:  "",      //‚	130	0082	BREAK PERMITTED HERE
	131:  "",      //ƒ	131	0083	NO BREAK HERE
	132:  "",      //„	132	0084	INDEX
	133:  "",      //…	133	0085	NEXT LINE (NEL)
	134:  "",      //†	134	0086	START OF SELECTED AREA
	135:  "",      //‡	135	0087	END OF SELECTED AREA
	136:  "",      //ˆ	136	0088	CHARACTER TABULATION SET
	137:  "",      //‰	137	0089	CHARACTER TABULATION WITH JUSTIFICATION
	138:  "",      //Š	138	008A	LINE TABULATION SET
	139:  "",      //‹	139	008B	PARTIAL LINE FORWARD
	140:  "",      //Œ	140	008C	PARTIAL LINE BACKWARD
	141:  "",      //	141	008D	REVERSE LINE FEED
	142:  "",      //Ž	142	008E	SINGLE SHIFT TWO
	143:  "",      //	143	008F	SINGLE SHIFT THREE
	144:  "",      //	144	0090	DEVICE CONTROL STRING
	145:  "",      //‘	145	0091	PRIVATE USE ONE
	146:  "",      //’	146	0092	PRIVATE USE TWO
	147:  "",      //“	147	0093	SET TRANSMIT STATE
	148:  "",      //”	148	0094	CANCEL CHARACTER
	149:  "",      //•	149	0095	MESSAGE WAITING
	150:  "",      //–	150	0096	START OF GUARDED AREA
	151:  "",      //—	151	0097	END OF GUARDED AREA
	152:  "",      //˜	152	0098	START OF STRING
	153:  "",      //™	153	0099	CONTROL
	154:  "",      //š	154	009A	SINGLE CHARACTER INTRODUCER
	155:  "",      //›	155	009B	CONTROL SEQUENCE INTRODUCER
	156:  "",      //œ	156	009C	STRING TERMINATOR
	157:  "",      //	157	009D	OPERATING SYSTEM COMMAND
	158:  "",      //ž	158	009E	PRIVACY MESSAGE
	159:  "",      //Ÿ	159	009F	APPLICATION PROGRAM COMMAND
	160:  " ",     //	160	00A0	&nbsp;	NO-BREAK SPACE
	161:  "",      //¡	161	00A1	&iexcl;	INVERTED EXCLAMATION MARK
	162:  "",      //¢	162	00A2	&cent;	CENT SIGN
	163:  "",      //£	163	00A3	&pound;	POUND SIGN
	164:  "",      //¤	164	00A4	&curren;	CURRENCY SIGN
	165:  "",      //¥	165	00A5	&yen;	YEN SIGN
	166:  "",      //¦	166	00A6	&brvbar;	BROKEN BAR
	167:  "",      //§	167	00A7	&sect;	SECTION SIGN
	168:  "",      //¨	168	00A8	&uml;	DIAERESIS
	169:  "",      //©	169	00A9	&copy;	COPYRIGHT SIGN
	170:  "",      //ª	170	00AA	&ordf;	FEMININE ORDINAL INDICATOR
	171:  "",      //«	171	00AB	&laquo;	LEFT-POINTING DOUBLE ANGLE QUOTATION MARK
	172:  "",      //¬	172	00AC	&not;	NOT SIGN
	173:  "",      //­	173	00AD	&shy;	SOFT HYPHEN
	174:  "",      //®	174	00AE	&reg;	REGISTERED SIGN
	175:  "",      //¯	175	00AF	&macr;	MACRON
	176:  "",      //°	176	00B0	&deg;	DEGREE SIGN
	177:  "",      //±	177	00B1	&plusmn;	PLUS-MINUS SIGN
	178:  "",      //²	178	00B2	&sup2;	SUPERSCRIPT TWO
	179:  "",      //³	179	00B3	&sup3;	SUPERSCRIPT THREE
	180:  "",      //´	180	00B4	&acute;	ACUTE ACCENT
	181:  "",      //µ	181	00B5	&micro;	MICRO SIGN
	182:  "",      //¶	182	00B6	&para;	PILCROW SIGN
	183:  "",      //·	183	00B7	&middot;	MIDDLE DOT
	184:  "",      //¸	184	00B8	&cedil;	CEDILLA
	185:  "",      //¹	185	00B9	&sup1;	SUPERSCRIPT ONE
	186:  "",      //º	186	00BA	&ordm;	MASCULINE ORDINAL INDICATOR
	187:  "",      //»	187	00BB	&raquo;	RIGHT-POINTING DOUBLE ANGLE QUOTATION MARK
	188:  "",      //¼	188	00BC	&frac14;	VULGAR FRACTION ONE QUARTER
	189:  "",      //½	189	00BD	&frac12;	VULGAR FRACTION ONE HALF
	190:  "",      //¾	190	00BE	&frac34;	VULGAR FRACTION THREE QUARTERS
	191:  "",      //¿	191	00BF	&iquest;	INVERTED QUESTION MARK
	192:  "A",     // À	192	00C0	&Agrave;	LATIN CAPITAL LETTER A WITH GRAVE
	193:  "A",     // Á	193
	194:  "A",     // Â	194
	195:  "A",     // Ã	195
	196:  "Ä",     // Ä	196
	197:  "A",     // Å	197
	198:  "A",     // Æ	198
	199:  "C",     //Ç	199	00C7	&Ccedil;	LATIN CAPITAL LETTER C WITH CEDILLA
	200:  "E",     //È	200	00C8	&Egrave;	LATIN CAPITAL LETTER E WITH GRAVE
	201:  "E",     //É	201	00C9	&Eacute;	LATIN CAPITAL LETTER E WITH ACUTE
	202:  "E",     //Ê	202	00CA	&Ecirc;	LATIN CAPITAL LETTER E WITH CIRCUMFLEX
	203:  "E",     //Ë	203	00CB	&Euml;	LATIN CAPITAL LETTER E WITH DIAERESIS
	204:  "I",     //Ì	204	00CC	&Igrave;	LATIN CAPITAL LETTER I WITH GRAVE
	205:  "I",     //Í	205	00CD	&Iacute;	LATIN CAPITAL LETTER I WITH ACUTE
	206:  "I",     //Î	206	00CE	&Icirc;	LATIN CAPITAL LETTER I WITH CIRCUMFLEX
	207:  "I",     //Ï	207	00CF	&Iuml;	LATIN CAPITAL LETTER I WITH DIAERESIS
	208:  "D",     //Ð	208	00D0	&ETH;	LATIN CAPITAL LETTER ETH
	209:  "N",     //Ñ	209	00D1	&Ntilde;	LATIN CAPITAL LETTER N WITH TILDE
	210:  "O",     //Ò	210	00D2	&Ograve;	LATIN CAPITAL LETTER O WITH GRAVE
	211:  "O",     //Ó	211	00D3	&Oacute;	LATIN CAPITAL LETTER O WITH ACUTE
	212:  "O",     //Ô	212	00D4	&Ocirc;	LATIN CAPITAL LETTER O WITH CIRCUMFLEX
	213:  "O",     //Õ	213	00D5	&Otilde;	LATIN CAPITAL LETTER O WITH TILDE
	214:  "Ö",     //Ö	214	00D6	&Ouml;	LATIN CAPITAL LETTER O WITH DIAERESIS
	215:  "",      //×	215	00D7	&times;	MULTIPLICATION SIGN
	216:  "O",     //Ø	216	00D8	&Oslash;	LATIN CAPITAL LETTER O WITH STROKE
	217:  "U",     //Ù	217	00D9	&Ugrave;	LATIN CAPITAL LETTER U WITH GRAVE
	218:  "U",     //Ú	218	00DA	&Uacute;	LATIN CAPITAL LETTER U WITH ACUTE
	219:  "U",     //Û	219	00DB	&Ucirc;	LATIN CAPITAL LETTER U WITH CIRCUMFLEX
	220:  "Ü",     //Ü	220	00DC	&Uuml;	LATIN CAPITAL LETTER U WITH DIAERESIS
	221:  "Y",     // Ý	221	00DD	&Yacute;	LATIN CAPITAL LETTER Y WITH ACUTE
	222:  "",      // Þ	222	00DE	&THORN;	LATIN CAPITAL LETTER THORN
	223:  "ß",     // ß
	224:  "a",     //à	224	00E0	&agrave;	LATIN SMALL LETTER A WITH GRAVE
	225:  "a",     //á	225	00E1	&aacute;	LATIN SMALL LETTER A WITH ACUTE
	226:  "a",     //â	226	00E2	&acirc;	LATIN SMALL LETTER A WITH CIRCUMFLEX
	227:  "a",     //ã	227	00E3	&atilde;	LATIN SMALL LETTER A WITH TILDE
	228:  "ä",     // ä
	229:  "a",     // å	229	00E5	&aring;	LATIN SMALL LETTER A WITH RING ABOVE
	230:  "ae",    // æ	230	00E6	&aelig;	LATIN SMALL LETTER AE
	231:  "c",     //ç	231	00E7	&ccedil;	LATIN SMALL LETTER C WITH CEDILLA
	232:  "e",     // è
	233:  "e",     // é
	234:  "e",     // ê
	235:  "e",     // 	ë	235	00EB	&euml;	LATIN SMALL LETTER E WITH DIAERESIS
	236:  "i",     //ì	236	00EC	&igrave;	LATIN SMALL LETTER I WITH GRAVE
	237:  "i",     //í	237	00ED	&iacute;	LATIN SMALL LETTER I WITH ACUTE
	238:  "i",     //î	238	00EE	&icirc;	LATIN SMALL LETTER I WITH CIRCUMFLEX
	239:  "i",     //ï	239	00EF	&iuml;	LATIN SMALL LETTER I WITH DIAERESIS
	240:  "o",     //ï	239	00EF	&iuml;	LATIN SMALL LETTER I WITH DIAERESIS
	241:  "n",     //	ñ	241	00F1	&ntilde;	LATIN SMALL LETTER N WITH TILDE
	242:  "o",     //ò	242	00F2	&ograve;	LATIN SMALL LETTER O WITH GRAVE
	243:  "o",     //ó	243	00F3	&oacute;	LATIN SMALL LETTER O WITH ACUTE
	244:  "o",     //ô	244	00F4	&ocirc;	LATIN SMALL LETTER O WITH CIRCUMFLEX
	245:  "o",     //õ	245	00F5	&otilde;	LATIN SMALL LETTER O WITH TILDE
	246:  "ö",     // ö
	247:  "",      //÷	247	00F7	&divide;	DIVISION SIGN
	248:  "o",     //ø	248	00F8	&oslash;	LATIN SMALL LETTER O WITH STROKE
	249:  "u",     //ù	249	00F9	&ugrave;	LATIN SMALL LETTER U WITH GRAVE
	250:  "u",     //ú	250	00FA	&uacute;	LATIN SMALL LETTER U WITH ACUTE
	251:  "u",     //û	251	00FB	&ucirc;	LATIN SMALL LETTER U WITH CIRCUMFLEX
	252:  "ü",     // ü
	253:  "y",     // ý 253
	254:  "b",     //þ	254	00FE	&thorn;	LATIN SMALL LETTER THORN
	255:  "y",     //ÿ	255	00FF	&yuml;	LATIN SMALL LETTER Y WITH DIAERESIS
	256:  "A",     //Ā	256	0100	&Amacr;	LATIN CAPITAL LETTER A WITH MACRON
	257:  "a",     //ā	257	0101	&amacr;	LATIN SMALL LETTER A WITH MACRON
	258:  "A",     //Ă	258	0102	&Abreve;	LATIN CAPITAL LETTER A WITH BREVE
	259:  "a",     //ă	259	0103	&abreve;	LATIN SMALL LETTER A WITH BREVE
	260:  "A",     //Ą	260	0104	&Aogon;	LATIN CAPITAL LETTER A WITH OGONEK
	261:  "q",     //q	261	0105	&aogon;	LATIN SMALL LETTER A WITH OGONEK
	262:  "C",     //Ć	262	0106	&Cacute;	LATIN CAPITAL LETTER C WITH ACUTE
	263:  "c",     //ć	263	0107	&cacute;	LATIN SMALL LETTER C WITH ACUTE
	264:  "C",     //Ĉ	264	0108	&Ccirc;	LATIN CAPITAL LETTER C WITH CIRCUMFLEX
	265:  "c",     //ĉ	265	0109	&ccirc;	LATIN SMALL LETTER C WITH CIRCUMFLEX
	266:  "C",     //Ċ	266	010A	&Cdod;	LATIN CAPITAL LETTER C WITH DOT ABOVE
	267:  "c",     //ċ	267	010B	&cdot;	LATIN SMALL LETTER C WITH DOT ABOVE
	268:  "C",     //Č	268	010C	&Ccaron;	LATIN CAPITAL LETTER C WITH CARON
	269:  "c",     //č	269	010D	&ccaron;	LATIN SMALL LETTER C WITH CARON
	270:  "C",     //Ď	270	010E	&Dcaron;	LATIN CAPITAL LETTER D WITH CARON
	271:  "c",     //ď	271	010F	&dcaron;	LATIN SMALL LETTER D WITH CARON
	272:  "C",     //Đ	272	0110	&Dstrok;	LATIN CAPITAL LETTER D WITH STROKE
	273:  "c",     //đ	273	0111	&dstrok;	LATIN SMALL LETTER D WITH STROKE
	274:  "E",     //Ē	274	0112	&Emacr;	LATIN CAPITAL LETTER E WITH MACRON
	275:  "e",     //ē	275	0113	&emacr;	LATIN SMALL LETTER E WITH MACRON
	276:  "E",     //Ĕ	276	0114	 	LATIN CAPITAL LETTER E WITH BREVE
	277:  "e",     //ĕ	277	0115	 	LATIN SMALL LETTER E WITH BREVE
	278:  "E",     //Ė	278	0116	&Edot;	LATIN CAPITAL LETTER E WITH DOT ABOVE
	279:  "e",     //ė	279	0117	&edot;	LATIN SMALL LETTER E WITH DOT ABOVE
	280:  "E",     //Ę	280	0118	&Eogon;	LATIN CAPITAL LETTER E WITH OGONEK
	281:  "e",     //ę	281	0119	&eogon;	LATIN SMALL LETTER E WITH OGONEK
	282:  "E",     //Ě	282	011A	&Ecaron;	LATIN CAPITAL LETTER E WITH CARON
	283:  "e",     //ě	283	011B	&ecaron;	LATIN SMALL LETTER E WITH CARON
	284:  "G",     //Ĝ	284	011C	&Gcirc;	LATIN CAPITAL LETTER G WITH CIRCUMFLEX
	285:  "g",     //ĝ	285	011D	&gcirc;	LATIN SMALL LETTER G WITH CIRCUMFLEX
	286:  "G",     //Ğ	286	011E	&Gbreve;	LATIN CAPITAL LETTER G WITH BREVE
	287:  "g",     //ğ	287	011F	&gbreve;	LATIN SMALL LETTER G WITH BREVE
	288:  "G",     //Ġ	288	0120	&Gdot;	LATIN CAPITAL LETTER G WITH DOT ABOVE
	289:  "g",     //ġ	289	0121	&gdot;	LATIN SMALL LETTER G WITH DOT ABOVE
	290:  "G",     //Ģ	290	0122	&Gcedil;	LATIN CAPITAL LETTER G WITH CEDILLA
	291:  "g",     //ģ	291	0123	&gcedil;	LATIN SMALL LETTER G WITH CEDILLA
	292:  "H",     //Ĥ	292	0124	&Hcirc;	LATIN CAPITAL LETTER H WITH CIRCUMFLEX
	293:  "h",     //ĥ	293	0125	&hcirc;	LATIN SMALL LETTER H WITH CIRCUMFLEX
	294:  "H",     //Ħ	294	0126	&Hstrok;	LATIN CAPITAL LETTER H WITH STROKE
	295:  "h",     //ħ	295	0127	&hstrok;	LATIN SMALL LETTER H WITH STROKE
	296:  "i",     // Ĩ &Itilde;	LATIN CAPITAL LETTER I WITH TILDE
	297:  "i",     // ĩ	297	0129	&itilde;	LATIN SMALL LETTER I WITH TILDE
	298:  "i",     // Ī	298	012A	&Imacr;	LATIN CAPITAL LETTER I WITH MACRON
	299:  "i",     // ī	299	012B	&imacr;	LATIN SMALL LETTER I WITH MACRON
	300:  "i",     // Ĭ	300	012C	 	LATIN CAPITAL LETTER I WITH BREVE
	301:  "i",     // ĭ	301	012D	 	LATIN SMALL LETTER I WITH BREVE
	302:  "i",     // Į	302	012E	&Iogon;	LATIN CAPITAL LETTER I WITH OGONEK
	303:  "i",     // į	303	012F	&iogon;	LATIN SMALL LETTER I WITH OGONEK
	304:  "i",     // İ	304	0130	&Idot;	LATIN CAPITAL LETTER I WITH DOT ABOVE
	305:  "i",     // ı	305	0131	&inodot;	LATIN SMALL LETTER DOTLESS I
	306:  "ij",    // Ĳ	306	0132	&IJlog;	LATIN CAPITAL LIGATURE IJ
	307:  "ij",    // ĳ	307	0133	&ijlig;	LATIN SMALL LIGATURE IJ
	308:  "j",     //Ĵ	308	0134	&Jcirc;	LATIN CAPITAL LETTER J WITH CIRCUMFLEX
	309:  "j",     //ĵ	309	0135	&jcirc;	LATIN SMALL LETTER J WITH CIRCUMFLEX
	310:  "K",     // Ķ	310	0136	&Kcedil;	LATIN CAPITAL LETTER K WITH CEDILLA
	311:  "k",     //ķ	311	0137	&kcedli;	LATIN SMALL LETTER K WITH CEDILLA
	312:  "k",     //ĸ	312	0138	&kgreen;	LATIN SMALL LETTER KRA
	313:  "L",     // 	Ĺ	313	0139	&Lacute;	LATIN CAPITAL LETTER L WITH ACUTE
	314:  "l",     // 	ĺ	314	013A	&lacute;	LATIN SMALL LETTER L WITH ACUTE
	315:  "L",     //Ļ	315	013B	&Lcedil;	LATIN CAPITAL LETTER L WITH CEDILLA
	316:  "l",     //ļ	316	013C	&lcedil;	LATIN SMALL LETTER L WITH CEDILLA
	317:  "L",     //Ľ	317	013D	&Lcaron;	LATIN CAPITAL LETTER L WITH CARON
	318:  "l",     //ľ	318	013E	&lcaron;	LATIN SMALL LETTER L WITH CARON
	319:  "L",     //Ŀ	319	013F	&Lmodot;	LATIN CAPITAL LETTER L WITH MIDDLE DOT
	320:  "l",     //ŀ	320	0140	&lmidot;	LATIN SMALL LETTER L WITH MIDDLE DOT
	321:  "L",     //Ł	321	0141	&Lstrok;	LATIN CAPITAL LETTER L WITH STROKE
	322:  "l",     //ł	322	0142	&lstrok;	LATIN SMALL LETTER L WITH STROKE
	323:  "N",     //Ń	323	0143	&Nacute;	LATIN CAPITAL LETTER N WITH ACUTE
	324:  "n",     //ń	324	0144	&nacute;	LATIN SMALL LETTER N WITH ACUTE
	325:  "N",     //Ņ	325	0145	&Ncedil;	LATIN CAPITAL LETTER N WITH CEDILLA
	326:  "n",     //ņ	326	0146	&ncedil;	LATIN SMALL LETTER N WITH CEDILLA
	327:  "N",     //Ň	327	0147	&Ncaron;	LATIN CAPITAL LETTER N WITH CARON
	328:  "n",     //ň	328	0148	&ncaron;	LATIN SMALL LETTER N WITH CARON
	329:  "n",     //ŉ	329	0149	&napos;	LATIN SMALL LETTER N PRECEDED BY APOSTROPHE
	330:  "N",     //Ŋ	330	014A	&ENG;	LATIN CAPITAL LETTER ENG
	331:  "n",     //ŋ	331	014B	&eng;	LATIN SMALL LETTER ENG
	332:  "O",     //Ō	332	014C	&Omacr;	LATIN CAPITAL LETTER O WITH MACRON
	333:  "o",     //ō	333	014D	&omacr;	LATIN SMALL LETTER O WITH MACRON
	334:  "O",     //Ŏ	334	014E	 	LATIN CAPITAL LETTER O WITH BREVE
	335:  "o",     //ŏ	335	014F	 	LATIN SMALL LETTER O WITH BREVE
	336:  "O",     //Ő	336	0150	&Odblac;	LATIN CAPITAL LETTER O WITH DOUBLE ACUTE
	337:  "ö",     //ő	337	0151	&odblac;	LATIN SMALL LETTER O WITH DOUBLE ACUTE
	338:  "Oe",    //Œ	338	0152	&OElig;	LATIN CAPITAL LIGATURE OE
	339:  "oe",    //œ	339	0153	&oelig;	LATIN SMALL LIGATURE OE
	340:  "R",     //Ŕ	340	0154	&Racute;	LATIN CAPITAL LETTER R WITH ACUTE
	341:  "r",     //ŕ	341	0155	&racute;	LATIN SMALL LETTER R WITH ACUTE
	342:  "R",     //Ŗ	342	0156	&Rcedil;	LATIN CAPITAL LETTER R WITH CEDILLA
	343:  "r",     //ŗ	343	0157	&rcedil;	LATIN SMALL LETTER R WITH CEDILLA
	344:  "R",     //Ř	344	0158	&Rcaron;	LATIN CAPITAL LETTER R WITH CARON
	345:  "r",     //ř	345	0159	&rcaron;	LATIN SMALL LETTER R WITH CARON
	346:  "S",     //Ś	346	015A	&Sacute;	LATIN CAPITAL LETTER S WITH ACUTE
	347:  "s",     //ś	347	015B	&sacute;	LATIN SMALL LETTER S WITH ACUTE
	348:  "S",     //Ŝ	348	015C	&Scirc;	LATIN CAPITAL LETTER S WITH CIRCUMFLEX
	349:  "s",     //ŝ	349	015D	&scirc;	LATIN SMALL LETTER S WITH CIRCUMFLEX
	350:  "S",     //Ş	350	015E	&Scedil;	LATIN CAPITAL LETTER S WITH CEDILLA
	351:  "s",     //ş	351	015F	&scedil;	LATIN SMALL LETTER S WITH CEDILLA
	352:  "S",     //Š	352	0160	&Scaron;	LATIN CAPITAL LETTER S WITH CARON
	353:  "s",     //š	353	0161	&scaron;	LATIN SMALL LETTER S WITH CARON
	354:  "T",     //Ţ	354	0162	&Tcedil;	LATIN CAPITAL LETTER T WITH CEDILLA
	355:  "t",     //ţ	355	0163	&tcedil;	LATIN SMALL LETTER T WITH CEDILLA
	356:  "T",     //Ť	356	0164	&Tcaron;	LATIN CAPITAL LETTER T WITH CARON
	357:  "t",     //ť	357	0165	&tcaron;	LATIN SMALL LETTER T WITH CARON
	358:  "T",     //Ŧ	358	0166	&Tstrok;	LATIN CAPITAL LETTER T WITH STROKE
	359:  "t",     //ŧ	359	0167	&tstrok;	LATIN SMALL LETTER T WITH STROKE
	360:  "U",     //Ũ	360	0168	&Utilde;	LATIN CAPITAL LETTER U WITH TILDE
	361:  "u",     //ũ	361	0169	&utilde;	LATIN SMALL LETTER U WITH TILDE
	362:  "U",     //Ū	362	016A	&Umacr;	LATIN CAPITAL LETTER U WITH MACRON
	363:  "u",     //ū	363	016B	&umacr;	LATIN SMALL LETTER U WITH MACRON
	364:  "U",     //Ŭ	364	016C	&Ubreve;	LATIN CAPITAL LETTER U WITH BREVE
	365:  "u",     //ŭ	365	016D	&ubreve;	LATIN SMALL LETTER U WITH BREVE
	366:  "U",     //Ů	366	016E	&Uring;	LATIN CAPITAL LETTER U WITH RING ABOVE
	367:  "u",     //ů	367	016F	&uring;	LATIN SMALL LETTER U WITH RING ABOVE
	368:  "U",     //Ű	368	0170	&Udblac;	LATIN CAPITAL LETTER U WITH DOUBLE ACUTE
	369:  "U",     //ű	369	0171	&udblac;	LATIN SMALL LETTER U WITH DOUBLE ACUTE
	370:  "U",     //Ų	370	0172	&Uogon;	LATIN CAPITAL LETTER U WITH OGONEK
	371:  "u",     //ų	371	0173	&uogon;	LATIN SMALL LETTER U WITH OGONEK
	372:  "W",     //Ŵ	372	0174	&Wcirc;	LATIN CAPITAL LETTER W WITH CIRCUMFLEX
	373:  "w",     //ŵ	373	0175	&wcirc;	LATIN SMALL LETTER W WITH CIRCUMFLEX
	374:  "Y",     //Ŷ	374	0176	&Ycirc;	LATIN CAPITAL LETTER Y WITH CIRCUMFLEX
	375:  "y",     //ŷ	375	0177	&ycirc;	LATIN SMALL LETTER Y WITH CIRCUMFLEX
	376:  "Y",     //Ÿ	376	0178	&Yuml;	LATIN CAPITAL LETTER Y WITH DIAERESIS
	377:  "Z",     //Ź	377	0179	&Zacute;	LATIN CAPITAL LETTER Z WITH ACUTE
	378:  "z",     //ź	378	017A	&zacute;	LATIN SMALL LETTER Z WITH ACUTE
	379:  "Z",     //Ż	379	017B	&Zdot;	LATIN CAPITAL LETTER Z WITH DOT ABOVE
	380:  "z",     // ż
	381:  "Z",     //Ž	381	017D	&Zcaron;	LATIN CAPITAL LETTER Z WITH CARON
	382:  "z",     //ž	382	017E	&zcaron;	LATIN SMALL LETTER Z WITH CARON
	383:  "s",     //ſ	383	017F	 	LATIN SMALL LETTER LONG S
	384:  "b",     //ƀ	384	0180	 	LATIN SMALL LETTER B WITH STROKE
	385:  "B",     //Ɓ	385	0181	 	LATIN CAPITAL LETTER B WITH HOOK
	386:  "B",     //Ƃ	386	0182	 	LATIN CAPITAL LETTER B WITH TOPBAR
	387:  "B",     //	ƃ	387	0183	 	LATIN SMALL LETTER B WITH TOPBAR
	388:  "b",     //Ƅ	388	0184	 	LATIN CAPITAL LETTER TONE SIX
	389:  "b",     //ƅ	389	0185	 	LATIN SMALL LETTER TONE SIX
	390:  "O",     //Ɔ	390	0186	 	LATIN CAPITAL LETTER OPEN O
	391:  "C",     //Ƈ	391	0187	 	LATIN CAPITAL LETTER C WITH HOOK
	392:  "c",     //ƈ	392	0188	 	LATIN SMALL LETTER C WITH HOOK
	393:  "D",     //Ɖ	393	0189	 	LATIN CAPITAL LETTER AFRICAN D
	394:  "D",     //Ɗ	394	018A	 	LATIN CAPITAL LETTER D WITH HOOK
	395:  "D",     //Ƌ	395	018B	 	LATIN CAPITAL LETTER D WITH TOPBAR
	396:  "D",     //ƌ	396	018C	 	LATIN SMALL LETTER D WITH TOPBAR
	397:  "d",     //ƍ	397	018D	 	LATIN SMALL LETTER TURNED DELTA
	398:  "E",     //Ǝ	398	018E	 	LATIN CAPITAL LETTER REVERSED E
	399:  "sch",   //Ə	399	018F	 	LATIN CAPITAL LETTER SCHWA
	400:  "e",     //Ɛ	400	0190	 	LATIN CAPITAL LETTER OPEN E
	401:  "F",     //Ƒ	401	0191	 	LATIN CAPITAL LETTER F WITH HOOK
	402:  "Ff",    //Ƒƒ	402	0192	&fnof;	LATIN SMALL LETTER F WITH HOOK
	403:  "Fg",    //ƑƓ	403	0193	 	LATIN CAPITAL LETTER G WITH HOOK
	404:  "Fy",    //ƑƔ	404	0194	 	LATIN CAPITAL LETTER GAMMA
	405:  "I",     //Ƒƕ	405	0195	 	LATIN SMALL LETTER HV
	406:  "I",     //ƑƖ	406	0196	 	LATIN CAPITAL LETTER IOTA
	407:  "I",     //ƑƗ	407	0197	 	LATIN CAPITAL LETTER I WITH STROKE
	408:  "K",     //ƑƘ	408	0198	 	LATIN CAPITAL LETTER K WITH HOOK
	409:  "k",     //Ƒƙ	409	0199	 	LATIN SMALL LETTER K WITH HOOK
	410:  "l",     //Ƒƚ	410	019A	 	LATIN SMALL LETTER L WITH BAR
	411:  "l",     //Ƒƛ	411	019B	 	LATIN SMALL LETTER LAMBDA WITH STROKE
	412:  "M",     //ƑƜ	412	019C	 	LATIN CAPITAL LETTER TURNED M
	413:  "N",     //ƑƝ	413	019D	 	LATIN CAPITAL LETTER N WITH LEFT HOOK
	414:  "n",     //Ƒƞ	414	019E	 	LATIN SMALL LETTER N WITH LONG RIGHT LEG
	415:  "O",     //ƑƟ	415	019F	 	LATIN CAPITAL LETTER O WITH MIDDLE TILDE
	416:  "O",     //ƑƠ	416	01A0	 	LATIN CAPITAL LETTER O WITH HORN
	417:  "o",     //Ƒơ	417	01A1	 	LATIN SMALL LETTER O WITH HORN
	418:  "Oi",    //ƑƢ	418	01A2	 	LATIN CAPITAL LETTER OI
	419:  "oi",    //Ƒƣ	419	01A3	 	LATIN SMALL LETTER OI
	420:  "P",     //ƑƤ	420	01A4	 	LATIN CAPITAL LETTER P WITH HOOK
	421:  "p",     //Ƒƥ	421	01A5	 	LATIN SMALL LETTER P WITH HOOK
	422:  "yr",    //ƑƦ	422	01A6	 	LATIN LETTER YR
	423:  "",      //ƑƧ	423	01A7	 	LATIN CAPITAL LETTER TONE TWO
	424:  "",      //Ƒƨ	424	01A8	 	LATIN SMALL LETTER TONE TWO
	425:  "",      //ƑƩ	425	01A9	 	LATIN CAPITAL LETTER ESH
	426:  "",      //Ƒƪ	426	01AA	 	LATIN LETTER REVERSED ESH LOOP
	427:  "t",     //Ƒƫ	427	01AB	 	LATIN SMALL LETTER T WITH PALATAL HOOK
	428:  "T",     //ƑƬ	428	01AC	 	LATIN CAPITAL LETTER T WITH HOOK
	429:  "t",     //Ƒƭ	429	01AD	 	LATIN SMALL LETTER T WITH HOOK
	430:  "T",     //ƑƮ	430	01AE	 	LATIN CAPITAL LETTER T WITH RETROFLEX HOOK
	431:  "u",     //ƑƯ	431	01AF	 	LATIN CAPITAL LETTER U WITH HORN
	432:  "u",     //Ƒư	432	01B0	 	LATIN SMALL LETTER U WITH HORN
	433:  "Y",     //ƑƱ	433	01B1	 	LATIN CAPITAL LETTER UPSILON
	434:  "V",     //ƑƲ	434	01B2	 	LATIN CAPITAL LETTER V WITH HOOK
	435:  "Y",     //ƑƳ	435	01B3	 	LATIN CAPITAL LETTER Y WITH HOOK
	436:  "Y",     //Ƒƴ	436	01B4	 	LATIN SMALL LETTER Y WITH HOOK
	437:  "Z",     //ƑƵ	437	01B5	&imped;	LATIN CAPITAL LETTER Z WITH STROKE
	438:  "z",     //Ƒƶ	438	01B6	 	LATIN SMALL LETTER Z WITH STROKE
	439:  "E",     //ƑƷ	439	01B7	 	LATIN CAPITAL LETTER EZH
	440:  "E",     //ƑƸ	440	01B8	 	LATIN CAPITAL LETTER EZH REVERSED
	441:  "e",     //Ƒƹ	441	01B9	 	LATIN SMALL LETTER EZH REVERSED
	442:  "e",     //Ƒƺ	442	01BA	 	LATIN SMALL LETTER EZH WITH TAIL
	443:  "2",     //Ƒƻ	443	01BB	 	LATIN LETTER TWO WITH STROKE
	444:  "5",     //ƑƼ	444	01BC	 	LATIN CAPITAL LETTER TONE FIVE
	445:  "5",     //Ƒƽ	445	01BD	 	LATIN SMALL LETTER TONE FIVE
	446:  "",      //Ƒƾ	446	01BE	 	LATIN LETTER INVERTED GLOTTAL STOP WITH STROKE
	447:  "",      //Ƒƿ	447	01BF	 	LATIN LETTER WYNN
	448:  "",      //Ƒǀ	448	01C0	 	LATIN LETTER DENTAL CLICK
	449:  "",      //Ƒǁ	449	01C1	 	LATIN LETTER LATERAL CLICK
	450:  "",      //Ƒǂ	450	01C2	 	LATIN LETTER ALVEOLAR CLICK
	451:  "",      //Ƒǃ	451	01C3	 	LATIN LETTER RETROFLEX CLICK
	452:  "Dz",    //ƑǄ	452	01C4	 	LATIN CAPITAL LETTER DZ WITH CARON
	453:  "D",     //Ƒǅ	453	01C5	 	LATIN CAPITAL LETTER D WITH SMALL LETTER Z WITH CARON
	454:  "dz",    //Ƒǆ	454	01C6	 	LATIN SMALL LETTER DZ WITH CARON
	455:  "Lj",    //ƑǇ	455	01C7	 	LATIN CAPITAL LETTER LJ
	456:  "Lj",    //Ƒǈ	456	01C8	 	LATIN CAPITAL LETTER L WITH SMALL LETTER J
	457:  "lj",    //Ƒǉ	457	01C9	 	LATIN SMALL LETTER LJ
	458:  "Nj",    //ƑǊ	458	01CA	 	LATIN CAPITAL LETTER NJ
	459:  "Nj",    //Ƒǋ	459	01CB	 	LATIN CAPITAL LETTER N WITH SMALL LETTER J
	460:  "nj",    //Ƒǌ	460	01CC	 	LATIN SMALL LETTER NJ
	461:  "A",     //ƑǍ	461	01CD	 	LATIN CAPITAL LETTER A WITH CARON
	462:  "a",     //Ƒǎ	462	01CE	 	LATIN SMALL LETTER A WITH CARON
	463:  "I",     //ƑǏ	463	01CF	 	LATIN CAPITAL LETTER I WITH CARON
	464:  "i",     //Ƒǐ	464	01D0	 	LATIN SMALL LETTER I WITH CARON
	465:  "O",     //ƑǑ	465	01D1	 	LATIN CAPITAL LETTER O WITH CARON
	466:  "o",     //Ƒǒ	466	01D2	 	LATIN SMALL LETTER O WITH CARON
	467:  "U",     //ƑǓ	467	01D3	 	LATIN CAPITAL LETTER U WITH CARON
	468:  "u",     //Ƒǔ	468	01D4	 	LATIN SMALL LETTER U WITH CARON
	469:  "U",     //ƑǕ	469	01D5	 	LATIN CAPITAL LETTER U WITH DIAERESIS AND MACRON
	470:  "u",     //Ƒǖ	470	01D6	 	LATIN SMALL LETTER U WITH DIAERESIS AND MACRON
	471:  "U",     //ƑǗ	471	01D7	 	LATIN CAPITAL LETTER U WITH DIAERESIS AND ACUTE
	472:  "u",     //Ƒǘ	472	01D8	 	LATIN SMALL LETTER U WITH DIAERESIS AND ACUTE
	473:  "U",     //ƑǙ	473	01D9	 	LATIN CAPITAL LETTER U WITH DIAERESIS AND CARON
	474:  "u",     //Ƒǚ	474	01DA	 	LATIN SMALL LETTER U WITH DIAERESIS AND CARON
	475:  "U",     //ƑǛ	475	01DB	 	LATIN CAPITAL LETTER U WITH DIAERESIS AND GRAVE
	476:  "u",     //Ƒǜ	476	01DC	 	LATIN SMALL LETTER U WITH DIAERESIS AND GRAVE
	477:  "e",     //Ƒǝ	477	01DD	 	LATIN SMALL LETTER TURNED E
	478:  "A",     //ƑǞ	478	01DE	 	LATIN CAPITAL LETTER A WITH DIAERESIS AND MACRON
	479:  "a",     //Ƒǟ	479	01DF	 	LATIN SMALL LETTER A WITH DIAERESIS AND MACRON
	480:  "A",     //ƑǠ	480	01E0	 	LATIN CAPITAL LETTER A WITH DOT ABOVE AND MACRON
	481:  "a",     //Ƒǡ	481	01E1	 	LATIN SMALL LETTER A WITH DOT ABOVE AND MACRON
	482:  "Ae",    //ƑǢ	482	01E2	 	LATIN CAPITAL LETTER AE WITH MACRON
	483:  "ae",    //Ƒǣ	483	01E3	 	LATIN SMALL LETTER AE WITH MACRON
	484:  "G",     //ƑǤ	484	01E4	 	LATIN CAPITAL LETTER G WITH STROKE
	485:  "g",     //Ƒǥ	485	01E5	 	LATIN SMALL LETTER G WITH STROKE
	486:  "G",     //ƑǦ	486	01E6	 	LATIN CAPITAL LETTER G WITH CARON
	487:  "g",     //Ƒǧ	487	01E7	 	LATIN SMALL LETTER G WITH CARON
	488:  "K",     //ƑǨ	488	01E8	 	LATIN CAPITAL LETTER K WITH CARON
	489:  "k",     //Ƒǩ	489	01E9	 	LATIN SMALL LETTER K WITH CARON
	490:  "O",     // Ǫ	490	01EA	 	LATIN CAPITAL LETTER O WITH OGONEK
	491:  "o",     //ǫ	491	01EB	 	LATIN SMALL LETTER O WITH OGONEK
	492:  "O",     //Ǭ	492	01EC	 	LATIN CAPITAL LETTER O WITH OGONEK AND MACRON
	493:  "o",     //ǭ	493	01ED	 	LATIN SMALL LETTER O WITH OGONEK AND MACRON
	494:  "Ezh",   //Ǯ	494	01EE	 	LATIN CAPITAL LETTER EZH WITH CARON
	495:  "esz",   //ǯ	495	01EF	 	LATIN SMALL LETTER EZH WITH CARON
	496:  "j",     //ǰ	496	01F0	 	LATIN SMALL LETTER J WITH CARON
	497:  "Dz",    //Ǳ	497	01F1	 	LATIN CAPITAL LETTER DZ
	498:  "Dz",    //ǲ	498	01F2	 	LATIN CAPITAL LETTER D WITH SMALL LETTER Z
	499:  "dz",    //ǳ	499	01F3	 	LATIN SMALL LETTER DZ
	500:  "G",     //Ǵ	500	01F4	 	LATIN CAPITAL LETTER G WITH ACUTE
	501:  "g",     //ǵ	501	01F5	&gacute;	LATIN SMALL LETTER G WITH ACUTE
	502:  "Hwair", //Ƕ	502	01F6	 	LATIN CAPITAL LETTER HWAIR
	503:  "Wynn",  //Ƿ	503	01F7	 	LATIN CAPITAL LETTER WYNN
	504:  "N",     //Ǹ	504	01F8	 	LATIN CAPITAL LETTER N WITH GRAVE
	505:  "n",     //ǹ	505	01F9	 	LATIN SMALL LETTER N WITH GRAVE
	506:  "A",     //Ǻ	506	01FA	 	LATIN CAPITAL LETTER A WITH RING ABOVE AND ACUTE
	507:  "a",     //ǻ	507	01FB	 	LATIN SMALL LETTER A WITH RING ABOVE AND ACUTE
	508:  "Ae",    //Ǽ	508	01FC	 	LATIN CAPITAL LETTER AE WITH ACUTE
	509:  "ae",    //ǽ	509	01FD	 	LATIN SMALL LETTER AE WITH ACUTE
	510:  "O",     //Ǿ	510	01FE	 	LATIN CAPITAL LETTER O WITH STROKE AND ACUTE
	511:  "o",     //ǿ	511	01FF	 	LATIN SMALL LETTER O WITH STROKE AND ACUTE
	512:  "A",     //Ȁ	512	0200	 	LATIN CAPITAL LETTER A WITH DOUBLE GRAVE
	513:  "a",     //ȁ	513	0201	 	LATIN SMALL LETTER A WITH DOUBLE GRAVE
	514:  "A",     //Ȃ	514	0202	 	LATIN CAPITAL LETTER A WITH INVERTED BREVE
	515:  "a",     //ȃ	515	0203	 	LATIN SMALL LETTER A WITH INVERTED BREVE
	516:  "E",     //Ȅ	516	0204	 	LATIN CAPITAL LETTER E WITH DOUBLE GRAVE
	517:  "e",     //ȅ	517	0205	 	LATIN SMALL LETTER E WITH DOUBLE GRAVE
	518:  "E",     //Ȇ	518	0206	 	LATIN CAPITAL LETTER E WITH INVERTED BREVE
	519:  "e",     //ȇ	519	0207	 	LATIN SMALL LETTER E WITH INVERTED BREVE
	520:  "E",     //Ȉ	520	0208	 	LATIN CAPITAL LETTER I WITH DOUBLE GRAVE
	521:  "i",     //ȉ	521	0209	 	LATIN SMALL LETTER I WITH DOUBLE GRAVE
	522:  "I",     //Ȋ	522	020A	 	LATIN CAPITAL LETTER I WITH INVERTED BREVE
	523:  "i",     //ȋ	523	020B	 	LATIN SMALL LETTER I WITH INVERTED BREVE
	524:  "O",     //Ȍ	524	020C	 	LATIN CAPITAL LETTER O WITH DOUBLE GRAVE
	525:  "o",     //ȍ	525	020D	 	LATIN SMALL LETTER O WITH DOUBLE GRAVE
	526:  "O",     //Ȏ	526	020E	 	LATIN CAPITAL LETTER O WITH INVERTED BREVE
	527:  "o",     //ȏ	527	020F	 	LATIN SMALL LETTER O WITH INVERTED BREVE
	528:  "R",     //Ȑ	528	0210	 	LATIN CAPITAL LETTER R WITH DOUBLE GRAVE
	529:  "r",     //ȑ	529	0211	 	LATIN SMALL LETTER R WITH DOUBLE GRAVE
	530:  "R",     //Ȓ	530	0212	 	LATIN CAPITAL LETTER R WITH INVERTED BREVE
	531:  "r",     //ȓ	531	0213	 	LATIN SMALL LETTER R WITH INVERTED BREVE
	532:  "U",     //Ȕ	532	0214	 	LATIN CAPITAL LETTER U WITH DOUBLE GRAVE
	533:  "u",     //ȕ	533	0215	 	LATIN SMALL LETTER U WITH DOUBLE GRAVE
	534:  "U",     //Ȗ	534	0216	 	LATIN CAPITAL LETTER U WITH INVERTED BREVE
	535:  "u",     //ȗ	535	0217	 	LATIN SMALL LETTER U WITH INVERTED BREVE
	536:  "S",     //Ș	536	0218	 	LATIN CAPITAL LETTER S WITH COMMA BELOW
	537:  "s",     //ș	537	0219	 	LATIN SMALL LETTER S WITH COMMA BELOW
	538:  "T",     //Ț	538	021A	 	LATIN CAPITAL LETTER T WITH COMMA BELOW
	539:  "t",     //ț	539	021B	 	LATIN SMALL LETTER T WITH COMMA BELOW
	540:  "z",     //Ȝ	540	021C	 	LATIN CAPITAL LETTER YOGH
	541:  "z",     //ȝ	541	021D	 	LATIN SMALL LETTER YOGH
	542:  "H",     //Ȟ	542	021E	 	LATIN CAPITAL LETTER H WITH CARON
	543:  "h",     //ȟ	543	021F	 	LATIN SMALL LETTER H WITH CARON
	544:  "n",     //Ƞ	544	0220	 	LATIN CAPITAL LETTER N WITH LONG RIGHT LEG
	545:  "d",     //ȡ	545	0221	 	LATIN SMALL LETTER D WITH CURL
	546:  "Ou",    //Ȣ	546	0222	 	LATIN CAPITAL LETTER OU
	547:  "ou",    //ȣ	547	0223	 	LATIN SMALL LETTER OU
	548:  "Z",     //Ȥ	548	0224	 	LATIN CAPITAL LETTER Z WITH HOOK
	549:  "z",     //ȥ	549	0225	 	LATIN SMALL LETTER Z WITH HOOK
	550:  "A",     //Ȧ	550	0226	 	LATIN CAPITAL LETTER A WITH DOT ABOVE
	551:  "a",     //ȧ	551	0227	 	LATIN SMALL LETTER A WITH DOT ABOVE
	552:  "E",     //Ȩ	552	0228	 	LATIN CAPITAL LETTER E WITH CEDILLA
	553:  "e",     //ȩ	553	0229	 	LATIN SMALL LETTER E WITH CEDILLA
	554:  "O",     //Ȫ	554	022A	 	LATIN CAPITAL LETTER O WITH DIAERESIS AND MACRON
	555:  "o",     //ȫ	555	022B	 	LATIN SMALL LETTER O WITH DIAERESIS AND MACRON
	556:  "O",     //Ȭ	556	022C	 	LATIN CAPITAL LETTER O WITH TILDE AND MACRON
	557:  "o",     //ȭ	557	022D	 	LATIN SMALL LETTER O WITH TILDE AND MACRON
	558:  "O",     //Ȯ	558	022E	 	LATIN CAPITAL LETTER O WITH DOT ABOVE
	559:  "o",     //ȯ	559	022F	 	LATIN SMALL LETTER O WITH DOT ABOVE
	560:  "O",     //Ȱ	560	0230	 	LATIN CAPITAL LETTER O WITH DOT ABOVE AND MACRON
	561:  "o",     //ȱ	561	0231	 	LATIN SMALL LETTER O WITH DOT ABOVE AND MACRON
	562:  "Y",     //Ȳ	562	0232	 	LATIN CAPITAL LETTER Y WITH MACRON
	563:  "y",     //ȳ	563	0233	 	LATIN SMALL LETTER Y WITH MACRON
	564:  "l",     //ȴ	564	0234	 	LATIN SMALL LETTER L WITH CURL
	565:  "n",     //ȵ	565	0235	 	LATIN SMALL LETTER N WITH CURL
	566:  "t",     //ȶ	566	0236	 	LATIN SMALL LETTER T WITH CURL
	567:  "j",     //ȷ	567	0237	&jmath;	LATIN SMALL LETTER DOTLESS J
	568:  "db",    //ȸ	568	0238	 	LATIN SMALL LETTER DB DIGRAPH
	569:  "qp",    //ȹ	569	0239	 	LATIN SMALL LETTER QP DIGRAPH
	570:  "A",     //Ⱥ	570	023A	 	LATIN CAPITAL LETTER A WITH STROKE
	571:  "C",     //Ȼ	571	023B	 	LATIN CAPITAL LETTER C WITH STROKE
	572:  "c",     //ȼ	572	023C	 	LATIN SMALL LETTER C WITH STROKE
	573:  "L",     //Ƚ	573	023D	 	LATIN CAPITAL LETTER L WITH BAR
	574:  "T",     //Ⱦ	574	023E	 	LATIN CAPITAL LETTER T WITH DIAGONAL STROKE
	575:  "s",     //ȿ	575	023F	 	LATIN SMALL LETTER S WITH SWASH TAIL
	576:  "z",     //ɀ	576	0240	 	LATIN SMALL LETTER Z WITH SWASH TAIL
	579:  "B",     //Ƀ	579	0243	 	LATIN CAPITAL LETTER B WITH STROKE
	580:  "U",     //Ʉ	580	0244	 	LATIN CAPITAL LETTER U BAR
	582:  "E",     //Ɇ	582	0246	 	LATIN CAPITAL LETTER E WITH STROKE
	583:  "e",     //ɇ	583	0247	 	LATIN SMALL LETTER E WITH STROKE
	584:  "J",     //Ɉ	584	0248	 	LATIN CAPITAL LETTER J WITH STROKE
	585:  "j",     //ɉ	585	0249	 	LATIN SMALL LETTER J WITH STROKE
	586:  "q",     //Ɋ	586	024A	 	LATIN CAPITAL LETTER SMALL Q WITH HOOK TAIL
	587:  "q",     //ɋ	587	024B	 	LATIN SMALL LETTER Q WITH HOOK TAIL
	588:  "R",     //Ɍ	588	024C	 	LATIN CAPITAL LETTER R WITH STROKE
	589:  "r",     //ɍ	589	024D	 	LATIN SMALL LETTER R WITH STROKE
	590:  "Y",     //Ɏ	590	024E	 	LATIN CAPITAL LETTER Y WITH STROKE
	591:  "y",     //ɏ	591	024F
	688:  "h",     //ʰ	688	02B0	 	MODIFIER LETTER SMALL H
	689:  "h",     //ʱ	689	02B1	 	MODIFIER LETTER SMALL H WITH HOOK
	690:  "j",     //ʲ	690	02B2	 	MODIFIER LETTER SMALL J
	691:  "r",     //ʳ	691	02B3	 	MODIFIER LETTER SMALL R
	692:  "r",     //ʴ	692	02B4	 	MODIFIER LETTER SMALL TURNED R
	693:  "r",     //ʵ	693	02B5	 	MODIFIER LETTER SMALL TURNED R WITH HOOK
	694:  "R",     //ʶ	694	02B6	 	MODIFIER LETTER SMALL CAPITAL INVERTED R
	695:  "w",     //ʷ	695	02B7	 	MODIFIER LETTER SMALL W
	696:  "y",     //ʸ	696	02B8	 	MODIFIER LETTER SMALL Y
	737:  "l",     //ˡ		737	02E1	 	MODIFIER LETTER SMALL L
	738:  "s",     //	ˢ	738	02E2	 	MODIFIER LETTER SMALL S
	739:  "x",     //ˣ		739	02E3	 	MODIFIER LETTER SMALL X
	768:  "o",     //ò	768	0300	 	GRAVE ACCENT
	769:  "o",     //ó	769	0301	 	ACUTE ACCENT
	770:  "o",     //ô	770	0302	 	CIRCUMFLEX ACCENT
	771:  "",      //õ	771	0303	 	TILDE
	772:  "o",     //ō	772	0304	 	MACRON
	773:  "o",     //o̅	773	0305	 	OVERLINE
	774:  "o",     //ŏ	774	0306	 	BREVE
	775:  "o",     //ȯ	775	0307	 	DOT ABOVE
	776:  "o",     //ö	776	0308	 	DIAERESIS
	777:  "o",     //ỏ	777	0309	 	HOOK ABOVE
	778:  "o",     //o̊	778	030A	 	RING ABOVE
	779:  "o",     //ő	779	030B	 	DOUBLE ACUTE ACCENT
	780:  "o",     //ǒ	780	030C	 	CARON
	781:  "o",     //o̍	781	030D	 	VERTICAL LINE ABOVE
	782:  "o",     //o̎	782	030E	 	DOUBLE VERTICAL LINE ABOVE
	783:  "o",     //ȍ	783	030F	 	DOUBLE GRAVE ACCENT
	784:  "o",     //o̐	784	0310	 	CANDRABINDU
	785:  "o",     //ȏ	785	0311	 	INVERTED BREVE
	786:  "o",     //o̒	786	0312	 	TURNED COMMA ABOVE
	787:  "o",     //o̓	787	0313	 	COMMA ABOVE
	788:  "o",     //o̔	788	0314	 	REVERSED COMMA ABOVE
	789:  "o",     //o̕	789	0315	 	COMMA ABOVE RIGHT
	790:  "o",     //o̖	790	0316	 	GRAVE ACCENT BELOW
	791:  "o",     //o̗	791	0317	 	ACUTE ACCENT BELOW
	792:  "o",     //o̘	792	0318	 	LEFT TACK BELOW
	793:  "o",     //o̙	793	0319	 	RIGHT TACK BELOW
	794:  "o",     //o̚	794	031A	 	LEFT ANGLE ABOVE
	795:  "o",     //ơ	795	031B	 	HORN
	796:  "o",     //o̜	796	031C	 	LEFT HALF RING BELOW
	797:  "o",     //o̝	797	031D	 	UP TACK BELOW
	798:  "o",     //o̞	798	031E	 	DOWN TACK BELOW
	799:  "o",     //o̟	799	031F	 	PLUS SIGN BELOW
	800:  "o",     //o̠	800	0320	 	MINUS SIGN BELOW
	801:  "o",     //o̡	801	0321	 	PALATALIZED HOOK BELOW
	802:  "o",     //o̢	802	0322	 	RETROFLEX HOOK BELOW
	803:  "o",     //ọ	803	0323	 	DOT BELOW
	804:  "o",     //o̤	804	0324	 	DIAERESIS BELOW
	805:  "o",     //o̥	805	0325	 	RING BELOW
	806:  "o",     //o̦	806	0326	 	COMMA BELOW
	807:  "o",     //o̧	807	0327	 	CEDILLA
	808:  "o",     //ǫ	808	0328	 	OGONEK
	809:  "o",     //o̩	809	0329	 	VERTICAL LINE BELOW
	810:  "o",     //o̪	810	032A	 	BRIDGE BELOW
	811:  "o",     //o̫	811	032B	 	INVERTED DOUBLE ARCH BELOW
	812:  "o",     //o̬	812	032C	 	CARON BELOW
	813:  "o",     //o̭	813	032D	 	CIRCUMFLEX ACCENT BELOW
	814:  "o",     //o̮	814	032E	 	BREVE BELOW
	815:  "o",     //o̯	815	032F	 	INVERTED BREVE BELOW
	816:  "o",     //o̰	816	0330	 	TILDE BELOW
	817:  "o",     //o̱	817	0331	 	MACRON BELOW
	818:  "o",     //o̲	818	0332	 	LOW LINE
	819:  "o",     //o̳	819	0333	 	DOUBLE LOW LINE
	820:  "o",     //o̴	820	0334	 	TILDE OVERLAY
	821:  "o",     //o̵	821	0335	 	SHORT STROKE OVERLAY
	822:  "o",     //o̶	822	0336	 	LONG STROKE OVERLAY
	823:  "o",     //o̷	823	0337	 	SHORT SOLIDUS OVERLAY
	824:  "o",     //o̸	824	0338	 	LONG SOLIDUS OVERLAY
	825:  "o",     //o̹	825	0339	 	RIGHT HALF RING BELOW
	826:  "o",     //o̺	826	033A	 	INVERTED BRIDGE BELOW
	827:  "o",     //o̻	827	033B	 	SQUARE BELOW
	828:  "o",     //o̼	828	033C	 	SEAGULL BELOW
	829:  "o",     //o̽	829	033D	 	X ABOVE
	830:  "o",     //o̾	830	033E	 	VERTICAL TILDE
	831:  "o",     //o̿	831	033F	 	DOUBLE OVERLINE
	832:  "o",     //ò	832	0340	 	GRAVE TONE MARK
	833:  "o",     //ó	833	0341	 	ACUTE TONE MARK
	834:  "o",     //o͂	834	0342	 	GREEK PERISPOMENI (combined with theta)
	835:  "o",     //o̓	835	0343	 	GREEK KORONIS (combined with theta)
	836:  "o",     //ö́	836	0344	 	GREEK DIALYTIKA TONOS (combined with theta)
	837:  "o",     //oͅ	837	0345	 	GREEK YPOGEGRAMMENI (combined with theta)
	838:  "o",     //o͆	838	0346	 	BRIDGE ABOVE
	839:  "o",     //o͇	839	0347	 	EQUALS SIGN BELOW
	840:  "o",     //o͈	840	0348	 	DOUBLE VERTICAL LINE BELOW
	841:  "o",     //o͉	841	0349	 	LEFT ANGLE BELOW
	842:  "o",     //	o͊	842	034A	 	NOT TILDE ABOVE
	843:  "o",     //	o͋	843	034B	 	HOMOTHETIC ABOVE
	844:  "o",     //o͌	844	034C	 	ALMOST EQUAL TO ABOVE
	845:  "o",     //o͍	845	034D	 	LEFT RIGHT ARROW BELOW
	846:  "o",     //o͎	846	034E	 	UPWARDS ARROW BELOW
	847:  "o",     //o͏	847	034F	 	GRAPHEME JOINER
	848:  "o",     //o͐	848	0350	 	RIGHT ARROWHEAD ABOVE
	849:  "o",     //o͑	849	0351	 	LEFT HALF RING ABOVE
	850:  "o",     //o͒	850	0352	 	FERMATA
	851:  "o",     //o͓	851	0353	 	X BELOW
	852:  "o",     //o͔	852	0354	 	LEFT ARROWHEAD BELOW
	853:  "o",     //o͕	853	0355	 	RIGHT ARROWHEAD BELOW
	854:  "o",     //o͖	854	0356	 	RIGHT ARROWHEAD AND UP ARROWHEAD BELOW
	855:  "o",     //o͗	855	0357	 	RIGHT HALF RING ABOVE
	856:  "o",     //o͘	856	0358	 	DOT ABOVE RIGHT
	857:  "o",     //o͙	857	0359	 	ASTERISK BELOW
	858:  "o",     //o͚	858	035A	 	DOUBLE RING BELOW
	859:  "o",     //	o͛	859	035B	 	ZIGZAG ABOVE
	860:  "o",     //	͜o	860	035C	 	DOUBLE BREVE BELOW
	861:  "o",     //	͝o	861	035D	 	DOUBLE BREVE
	862:  "o",     //	͞o	862	035E	 	DOUBLE MACRON
	863:  "o",     //	͟o	863	035F	 	DOUBLE MACRON BELOW
	864:  "o",     //	͠o	864	0360	 	DOUBLE TILDE
	865:  "o",     //	͡o	865	0361	 	DOUBLE INVERTED BREVE
	866:  "o",     //	͢o	866	0362	 	DOUBLE RIGHTWARDS ARROW BELOW
	867:  "o",     //	oͣ	867	0363	 	LATIN SMALL LETTER A
	868:  "e",     //oͤ	868	0364	 	LATIN SMALL LETTER E
	869:  "i",     //oͥ	869	0365	 	LATIN SMALL LETTER I
	870:  "o",     //oͦ	870	0366	 	LATIN SMALL LETTER O
	871:  "u",     //oͧ	871	0367	 	LATIN SMALL LETTER U
	872:  "c",     //oͨ	872	0368	 	LATIN SMALL LETTER C
	873:  "d",     //oͩ	873	0369	 	LATIN SMALL LETTER D
	874:  "h",     //oͪ	874	036A	 	LATIN SMALL LETTER H
	875:  "m",     //oͫ	875	036B	 	LATIN SMALL LETTER M
	876:  "r",     //oͬ	876	036C	 	LATIN SMALL LETTER R
	877:  "t",     //oͭ	877	036D	 	LATIN SMALL LETTER T
	878:  "v",     //oͮ	878	036E	 	LATIN SMALL LETTER V
	879:  "x",     //oͯ	879	036F	 	LATIN SMALL LETTER X
	1104: "e",     //	ѐ	1104	0450	 	CYRILLIC SMALL LETTER IE WITH GRAVE
	1105: "e",     //	ё	1105	0451	 	CYRILLIC SMALL LETTER IO
	1106: "dje",   //	ђ	1106	0452	 	CYRILLIC SMALL LETTER DJE
	1107: "gje",   //	ѓ	1107	0453	 	CYRILLIC SMALL LETTER GJE
	1108: "ie",    //	є	1108	0454	 	CYRILLIC SMALL LETTER UKRAINIAN IE
	1109: "dze",   //	ѕ	1109	0455	 	CYRILLIC SMALL LETTER DZE
	1110: "i",     //	і	1110	0456	 	CYRILLIC SMALL LETTER BYELORUSSIAN-UKRAINIAN I
	1111: "yi",    //	ї	1111	0457	 	CYRILLIC SMALL LETTER YI
	1112: "je",    //	ј	1112	0458	 	CYRILLIC SMALL LETTER JE
	1113: "lje",   //	љ	1113	0459	 	CYRILLIC SMALL LETTER LJE
	1114: "nje",   //	њ	1114	045A	 	CYRILLIC SMALL LETTER NJE
	1115: "tshe",  //	ћ	1115	045B	 	CYRILLIC SMALL LETTER TSHE
	1116: "kje",   //	ќ	1116	045C	 	CYRILLIC SMALL LETTER KJE
	1117: "i",     //	ѝ	1117	045D	 	CYRILLIC SMALL LETTER I WITH GRAVE
	7729: "k",     //    ḱ

	7840: "A", // Ạ
	7841: "a", //ạ
	7842: "A", // Ả
	7843: "a", //ả
	7844: "A", // Ấ
	7845: "a", // ấ 7845
	7846: "A", // Ầ
	7847: "a", //  ầ
	7848: "A", // Ẩ
	7849: "a", //ẩ
	7850: "A", //Ẫ
	7879: "e", // ệ 7879
	7891: "o", // ồ
	8357: "m", // ₥
}
