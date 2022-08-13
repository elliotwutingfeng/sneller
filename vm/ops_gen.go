package vm

// Code generated automatically; DO NOT EDIT

const (
	opret                    bcop = 0
	opjz                     bcop = 1
	oploadk                  bcop = 2
	opsavek                  bcop = 3
	opxchgk                  bcop = 4
	oploadb                  bcop = 5
	opsaveb                  bcop = 6
	oploadv                  bcop = 7
	opsavev                  bcop = 8
	oploadzerov              bcop = 9
	opsavezerov              bcop = 10
	oploadpermzerov          bcop = 11
	opsaveblendv             bcop = 12
	oploads                  bcop = 13
	opsaves                  bcop = 14
	oploadzeros              bcop = 15
	opsavezeros              bcop = 16
	opbroadcastimmk          bcop = 17
	opfalse                  bcop = 18
	opandk                   bcop = 19
	opork                    bcop = 20
	opandnotk                bcop = 21
	opnandk                  bcop = 22
	opxork                   bcop = 23
	opnotk                   bcop = 24
	opxnork                  bcop = 25
	opbroadcastimmf          bcop = 26
	opbroadcastimmi          bcop = 27
	opabsf                   bcop = 28
	opabsi                   bcop = 29
	opnegf                   bcop = 30
	opnegi                   bcop = 31
	opsignf                  bcop = 32
	opsigni                  bcop = 33
	opsquaref                bcop = 34
	opsquarei                bcop = 35
	opbitnoti                bcop = 36
	opbitcounti              bcop = 37
	oproundf                 bcop = 38
	oproundevenf             bcop = 39
	optruncf                 bcop = 40
	opfloorf                 bcop = 41
	opceilf                  bcop = 42
	opaddf                   bcop = 43
	opaddimmf                bcop = 44
	opaddi                   bcop = 45
	opaddimmi                bcop = 46
	opsubf                   bcop = 47
	opsubimmf                bcop = 48
	opsubi                   bcop = 49
	opsubimmi                bcop = 50
	oprsubf                  bcop = 51
	oprsubimmf               bcop = 52
	oprsubi                  bcop = 53
	oprsubimmi               bcop = 54
	opmulf                   bcop = 55
	opmulimmf                bcop = 56
	opmuli                   bcop = 57
	opmulimmi                bcop = 58
	opdivf                   bcop = 59
	opdivimmf                bcop = 60
	oprdivf                  bcop = 61
	oprdivimmf               bcop = 62
	opdivi                   bcop = 63
	opdivimmi                bcop = 64
	oprdivi                  bcop = 65
	oprdivimmi               bcop = 66
	opmodf                   bcop = 67
	opmodimmf                bcop = 68
	oprmodf                  bcop = 69
	oprmodimmf               bcop = 70
	opmodi                   bcop = 71
	opmodimmi                bcop = 72
	oprmodi                  bcop = 73
	oprmodimmi               bcop = 74
	opaddmulimmi             bcop = 75
	opminvaluef              bcop = 76
	opminvalueimmf           bcop = 77
	opmaxvaluef              bcop = 78
	opmaxvalueimmf           bcop = 79
	opminvaluei              bcop = 80
	opminvalueimmi           bcop = 81
	opmaxvaluei              bcop = 82
	opmaxvalueimmi           bcop = 83
	opandi                   bcop = 84
	opandimmi                bcop = 85
	opori                    bcop = 86
	oporimmi                 bcop = 87
	opxori                   bcop = 88
	opxorimmi                bcop = 89
	opslli                   bcop = 90
	opsllimmi                bcop = 91
	opsrai                   bcop = 92
	opsraimmi                bcop = 93
	opsrli                   bcop = 94
	opsrlimmi                bcop = 95
	opsqrtf                  bcop = 96
	opcbrtf                  bcop = 97
	opexpf                   bcop = 98
	opexp2f                  bcop = 99
	opexp10f                 bcop = 100
	opexpm1f                 bcop = 101
	oplnf                    bcop = 102
	opln1pf                  bcop = 103
	oplog2f                  bcop = 104
	oplog10f                 bcop = 105
	opsinf                   bcop = 106
	opcosf                   bcop = 107
	optanf                   bcop = 108
	opasinf                  bcop = 109
	opacosf                  bcop = 110
	opatanf                  bcop = 111
	opatan2f                 bcop = 112
	ophypotf                 bcop = 113
	oppowf                   bcop = 114
	opcvtktof64              bcop = 115
	opcvtktoi64              bcop = 116
	opcvti64tok              bcop = 117
	opcvti64tof64            bcop = 118
	opcvtf64toi64            bcop = 119
	opfproundu               bcop = 120
	opfproundd               bcop = 121
	opcvti64tostr            bcop = 122
	opcmpv                   bcop = 123
	opcmpvk                  bcop = 124
	opcmpvimmk               bcop = 125
	opcmpvi64                bcop = 126
	opcmpvimmi64             bcop = 127
	opcmpvf64                bcop = 128
	opcmpvimmf64             bcop = 129
	opcmpltstr               bcop = 130
	opcmplestr               bcop = 131
	opcmpgtstr               bcop = 132
	opcmpgestr               bcop = 133
	opcmpltk                 bcop = 134
	opcmpltimmk              bcop = 135
	opcmplek                 bcop = 136
	opcmpleimmk              bcop = 137
	opcmpgtk                 bcop = 138
	opcmpgtimmk              bcop = 139
	opcmpgek                 bcop = 140
	opcmpgeimmk              bcop = 141
	opcmpeqf                 bcop = 142
	opcmpeqi                 bcop = 143
	opcmpeqimmf              bcop = 144
	opcmpeqimmi              bcop = 145
	opcmpltf                 bcop = 146
	opcmplti                 bcop = 147
	opcmpltimmf              bcop = 148
	opcmpltimmi              bcop = 149
	opcmplef                 bcop = 150
	opcmplei                 bcop = 151
	opcmpleimmf              bcop = 152
	opcmpleimmi              bcop = 153
	opcmpgtf                 bcop = 154
	opcmpgti                 bcop = 155
	opcmpgtimmf              bcop = 156
	opcmpgtimmi              bcop = 157
	opcmpgef                 bcop = 158
	opcmpgei                 bcop = 159
	opcmpgeimmf              bcop = 160
	opcmpgeimmi              bcop = 161
	opisnanf                 bcop = 162
	opchecktag               bcop = 163
	opisnull                 bcop = 164
	opisnotnull              bcop = 165
	opistrue                 bcop = 166
	opisfalse                bcop = 167
	opeqslice                bcop = 168
	opequalv                 bcop = 169
	opeqv4mask               bcop = 170
	opeqv4maskplus           bcop = 171
	opeqv8                   bcop = 172
	opeqv8plus               bcop = 173
	opleneq                  bcop = 174
	opdateaddmonth           bcop = 175
	opdateaddmonthimm        bcop = 176
	opdateaddyear            bcop = 177
	opdatediffparam          bcop = 178
	opdatediffmonthyear      bcop = 179
	opdateextractmicrosecond bcop = 180
	opdateextractmillisecond bcop = 181
	opdateextractsecond      bcop = 182
	opdateextractminute      bcop = 183
	opdateextracthour        bcop = 184
	opdateextractday         bcop = 185
	opdateextractmonth       bcop = 186
	opdateextractyear        bcop = 187
	opdatetounixepoch        bcop = 188
	opdatetruncmillisecond   bcop = 189
	opdatetruncsecond        bcop = 190
	opdatetruncminute        bcop = 191
	opdatetrunchour          bcop = 192
	opdatetruncday           bcop = 193
	opdatetruncmonth         bcop = 194
	opdatetruncyear          bcop = 195
	opunboxts                bcop = 196
	opboxts                  bcop = 197
	optimelt                 bcop = 198
	optimegt                 bcop = 199
	opconsttm                bcop = 200
	optmextract              bcop = 201
	opwidthbucketf           bcop = 202
	opwidthbucketi           bcop = 203
	optimebucketts           bcop = 204
	opgeohash                bcop = 205
	opgeohashimm             bcop = 206
	opgeotilex               bcop = 207
	opgeotiley               bcop = 208
	opgeotilees              bcop = 209
	opgeotileesimm           bcop = 210
	opgeodistance            bcop = 211
	opconcatlenget1          bcop = 212
	opconcatlenget2          bcop = 213
	opconcatlenget3          bcop = 214
	opconcatlenget4          bcop = 215
	opconcatlenacc1          bcop = 216
	opconcatlenacc2          bcop = 217
	opconcatlenacc3          bcop = 218
	opconcatlenacc4          bcop = 219
	opallocstr               bcop = 220
	opappendstr              bcop = 221
	opfindsym                bcop = 222
	opfindsym2               bcop = 223
	opfindsym2rev            bcop = 224
	opfindsym3               bcop = 225
	opblendv                 bcop = 226
	opblendrevv              bcop = 227
	opblendnum               bcop = 228
	opblendnumrev            bcop = 229
	opblendslice             bcop = 230
	opblendslicerev          bcop = 231
	opunpack                 bcop = 232
	opunsymbolize            bcop = 233
	opunboxktoi64            bcop = 234
	optoint                  bcop = 235
	optof64                  bcop = 236
	opboxfloat               bcop = 237
	opboxint                 bcop = 238
	opboxmask                bcop = 239
	opboxmask2               bcop = 240
	opboxmask3               bcop = 241
	opboxstring              bcop = 242
	ophashvalue              bcop = 243
	ophashvalueplus          bcop = 244
	ophashmember             bcop = 245
	ophashlookup             bcop = 246
	opaggandk                bcop = 247
	opaggork                 bcop = 248
	opaggsumf                bcop = 249
	opaggsumi                bcop = 250
	opaggminf                bcop = 251
	opaggmini                bcop = 252
	opaggmaxf                bcop = 253
	opaggmaxi                bcop = 254
	opaggandi                bcop = 255
	opaggori                 bcop = 256
	opaggxori                bcop = 257
	opaggcount               bcop = 258
	opaggbucket              bcop = 259
	opaggslotandk            bcop = 260
	opaggslotork             bcop = 261
	opaggslotaddf            bcop = 262
	opaggslotaddi            bcop = 263
	opaggslotavgf            bcop = 264
	opaggslotavgi            bcop = 265
	opaggslotminf            bcop = 266
	opaggslotmini            bcop = 267
	opaggslotmaxf            bcop = 268
	opaggslotmaxi            bcop = 269
	opaggslotandi            bcop = 270
	opaggslotori             bcop = 271
	opaggslotxori            bcop = 272
	opaggslotcount           bcop = 273
	oplitref                 bcop = 274
	opsplit                  bcop = 275
	optuple                  bcop = 276
	opdupv                   bcop = 277
	opzerov                  bcop = 278
	opobjectsize             bcop = 279
	opCmpStrEqCs             bcop = 280
	opCmpStrEqCi             bcop = 281
	opCmpStrEqUTF8Ci         bcop = 282
	opSkip1charLeft          bcop = 283
	opSkip1charRight         bcop = 284
	opSkipNcharLeft          bcop = 285
	opSkipNcharRight         bcop = 286
	opTrimWsLeft             bcop = 287
	opTrimWsRight            bcop = 288
	opTrim4charLeft          bcop = 289
	opTrim4charRight         bcop = 290
	opContainsSubstrCs       bcop = 291
	opContainsSubstrCi       bcop = 292
	opContainsSuffixCs       bcop = 293
	opContainsSuffixCi       bcop = 294
	opContainsSuffixUTF8Ci   bcop = 295
	opContainsPrefixCs       bcop = 296
	opContainsPrefixCi       bcop = 297
	opContainsPrefixUTF8Ci   bcop = 298
	opLengthStr              bcop = 299
	opSubstr                 bcop = 300
	opSplitPart              bcop = 301
	opMatchpatCs             bcop = 302
	opMatchpatCi             bcop = 303
	opMatchpatUTF8Ci         bcop = 304
	opIsSubnetOfIP4          bcop = 305
	opDfaT6                  bcop = 306
	opDfaT7                  bcop = 307
	opDfaT8                  bcop = 308
	opDfaT6Z                 bcop = 309
	opDfaT7Z                 bcop = 310
	opDfaT8Z                 bcop = 311
	opDfaL                   bcop = 312
	opDfaLZ                  bcop = 313
	opslower                 bcop = 314
	opsupper                 bcop = 315
	opsadjustsize            bcop = 316
	optrap                   bcop = 317
	_maxbcop                      = 318
)
