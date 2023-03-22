package vm

// Code generated automatically; DO NOT EDIT

var generated = [_maxbcop]bcopinfo{
	opbroadcasti64:            {text: "broadcast.i64", out: bcargs[0:1] /* {bcS} */, in: bcargs[14:15] /* {bcImmI64} */},
	opabsi64:                  {text: "abs.i64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opnegi64:                  {text: "neg.i64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opsigni64:                 {text: "sign.i64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opsquarei64:               {text: "square.i64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opbitnoti64:               {text: "bitnot.i64", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opbitcounti64:             {text: "bitcount.i64", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opbitcounti64_v2:          {text: "bitcount.i64", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opaddi64:                  {text: "add.i64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opaddi64imm:               {text: "add.i64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opsubi64:                  {text: "sub.i64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opsubi64imm:               {text: "sub.i64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	oprsubi64imm:              {text: "rsub.i64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opmuli64:                  {text: "mul.i64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opmuli64imm:               {text: "mul.i64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opdivi64:                  {text: "div.i64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opdivi64imm:               {text: "div.i64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	oprdivi64imm:              {text: "rdiv.i64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opmodi64:                  {text: "mod.i64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opmodi64imm:               {text: "mod.i64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	oprmodi64imm:              {text: "rmod.i64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opaddmuli64imm:            {text: "addmul.i64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[12:16] /* {bcS, bcS, bcImmI64, bcK} */},
	opminvaluei64:             {text: "minvalue.i64", out: bcargs[0:1] /* {bcS} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opminvaluei64imm:          {text: "minvalue.i64@imm", out: bcargs[0:1] /* {bcS} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opmaxvaluei64:             {text: "maxvalue.i64", out: bcargs[0:1] /* {bcS} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opmaxvaluei64imm:          {text: "maxvalue.i64@imm", out: bcargs[0:1] /* {bcS} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opandi64:                  {text: "and.i64", out: bcargs[0:1] /* {bcS} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opandi64imm:               {text: "and.i64@imm", out: bcargs[0:1] /* {bcS} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opori64:                   {text: "or.i64", out: bcargs[0:1] /* {bcS} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opori64imm:                {text: "or.i64@imm", out: bcargs[0:1] /* {bcS} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opxori64:                  {text: "xor.i64", out: bcargs[0:1] /* {bcS} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opxori64imm:               {text: "xor.i64@imm", out: bcargs[0:1] /* {bcS} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opslli64:                  {text: "sll.i64", out: bcargs[0:1] /* {bcS} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opslli64imm:               {text: "sll.i64@imm", out: bcargs[0:1] /* {bcS} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opsrai64:                  {text: "sra.i64", out: bcargs[0:1] /* {bcS} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opsrai64imm:               {text: "sra.i64@imm", out: bcargs[0:1] /* {bcS} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opsrli64:                  {text: "srl.i64", out: bcargs[0:1] /* {bcS} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opsrli64imm:               {text: "srl.i64@imm", out: bcargs[0:1] /* {bcS} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opbroadcastf64:            {text: "broadcast.f64", out: bcargs[0:1] /* {bcS} */, in: bcargs[17:18] /* {bcImmF64} */},
	opabsf64:                  {text: "abs.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opnegf64:                  {text: "neg.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opsignf64:                 {text: "sign.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opsquaref64:               {text: "square.f64", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	oproundf64:                {text: "round.f64", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	oproundevenf64:            {text: "roundeven.f64", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	optruncf64:                {text: "trunc.f64", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opfloorf64:                {text: "floor.f64", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opceilf64:                 {text: "ceil.f64", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opaddf64:                  {text: "add.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opaddf64imm:               {text: "add.f64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[106:109] /* {bcS, bcImmF64, bcK} */},
	opsubf64:                  {text: "sub.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opsubf64imm:               {text: "sub.f64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[106:109] /* {bcS, bcImmF64, bcK} */},
	oprsubf64imm:              {text: "rsub.f64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[106:109] /* {bcS, bcImmF64, bcK} */},
	opmulf64:                  {text: "mul.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opmulf64imm:               {text: "mul.f64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[106:109] /* {bcS, bcImmF64, bcK} */},
	opdivf64:                  {text: "div.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opdivf64imm:               {text: "div.f64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[106:109] /* {bcS, bcImmF64, bcK} */},
	oprdivf64imm:              {text: "rdiv.f64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[106:109] /* {bcS, bcImmF64, bcK} */},
	opmodf64:                  {text: "mod.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opmodf64imm:               {text: "mod.f64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[106:109] /* {bcS, bcImmF64, bcK} */},
	oprmodf64imm:              {text: "rmod.f64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[106:109] /* {bcS, bcImmF64, bcK} */},
	opminvaluef64:             {text: "minvalue.f64", out: bcargs[0:1] /* {bcS} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opminvaluef64imm:          {text: "minvalue.f64@imm", out: bcargs[0:1] /* {bcS} */, in: bcargs[106:109] /* {bcS, bcImmF64, bcK} */},
	opmaxvaluef64:             {text: "maxvalue.f64", out: bcargs[0:1] /* {bcS} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opmaxvaluef64imm:          {text: "maxvalue.f64@imm", out: bcargs[0:1] /* {bcS} */, in: bcargs[106:109] /* {bcS, bcImmF64, bcK} */},
	opsqrtf64:                 {text: "sqrt.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opcbrtf64:                 {text: "cbrt.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opexpf64:                  {text: "exp.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opexp2f64:                 {text: "exp2.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opexp10f64:                {text: "exp10.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opexpm1f64:                {text: "expm1.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	oplnf64:                   {text: "ln.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opln1pf64:                 {text: "ln1p.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	oplog2f64:                 {text: "log2.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	oplog10f64:                {text: "log10.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opsinf64:                  {text: "sin.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opcosf64:                  {text: "cos.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	optanf64:                  {text: "tan.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opasinf64:                 {text: "asin.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opacosf64:                 {text: "acos.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opatanf64:                 {text: "atan.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opatan2f64:                {text: "atan2.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	ophypotf64:                {text: "hypot.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	oppowf64:                  {text: "pow.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	optrap:                    {text: "trap"},
	opret:                     {text: "ret"},
	opretk:                    {text: "ret.k", in: bcargs[4:5] /* {bcK} */},
	opretbk:                   {text: "ret.b.k", in: bcargs[44:46] /* {bcB, bcK} */},
	opretsk:                   {text: "ret.s.k", in: bcargs[3:5] /* {bcS, bcK} */},
	opretbhk:                  {text: "ret.b.h.k", in: bcargs[35:38] /* {bcB, bcH, bcK} */},
	opinit:                    {text: "init", out: bcargs[44:46] /* {bcB, bcK} */},
	opbroadcast0k:             {text: "broadcast0.k", out: bcargs[4:5] /* {bcK} */},
	opbroadcast1k:             {text: "broadcast1.k", out: bcargs[4:5] /* {bcK} */},
	opfalse:                   {text: "false.k", out: bcargs[10:12] /* {bcV, bcK} */},
	opnotk:                    {text: "not.k", out: bcargs[4:5] /* {bcK} */, in: bcargs[4:5] /* {bcK} */},
	opandk:                    {text: "and.k", out: bcargs[4:5] /* {bcK} */, in: bcargs[7:9] /* {bcK, bcK} */},
	opandnk:                   {text: "andn.k", out: bcargs[4:5] /* {bcK} */, in: bcargs[7:9] /* {bcK, bcK} */},
	opork:                     {text: "or.k", out: bcargs[4:5] /* {bcK} */, in: bcargs[7:9] /* {bcK, bcK} */},
	opxork:                    {text: "xor.k", out: bcargs[4:5] /* {bcK} */, in: bcargs[7:9] /* {bcK, bcK} */},
	opxnork:                   {text: "xnor.k", out: bcargs[4:5] /* {bcK} */, in: bcargs[7:9] /* {bcK, bcK} */},
	opcvtktof64:               {text: "cvt.ktof64", out: bcargs[0:1] /* {bcS} */, in: bcargs[4:5] /* {bcK} */},
	opcvtktoi64:               {text: "cvt.ktoi64", out: bcargs[0:1] /* {bcS} */, in: bcargs[4:5] /* {bcK} */},
	opcvti64tok:               {text: "cvt.i64tok", out: bcargs[4:5] /* {bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opcvtf64tok:               {text: "cvt.f64tok", out: bcargs[4:5] /* {bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opcvti64tof64:             {text: "cvt.i64tof64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opcvttruncf64toi64:        {text: "cvttrunc.f64toi64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opcvtfloorf64toi64:        {text: "cvtfloor.f64toi64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opcvtceilf64toi64:         {text: "cvtceil.f64toi64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opcvti64tostr:             {text: "cvt.i64tostr", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */, scratch: 20 * 16},
	opcmpv:                    {text: "cmpv", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[99:102] /* {bcV, bcV, bcK} */},
	opsortcmpvnf:              {text: "sortcmpv@nf", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[99:102] /* {bcV, bcV, bcK} */},
	opsortcmpvnl:              {text: "sortcmpv@nl", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[99:102] /* {bcV, bcV, bcK} */},
	opcmpvk:                   {text: "cmpv.k", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[72:75] /* {bcV, bcK, bcK} */},
	opcmpvkimm:                {text: "cmpv.k@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[53:56] /* {bcV, bcImmU16, bcK} */},
	opcmpvi64:                 {text: "cmpv.i64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[56:59] /* {bcV, bcS, bcK} */},
	opcmpvi64imm:              {text: "cmpv.i64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[75:78] /* {bcV, bcImmI64, bcK} */},
	opcmpvf64:                 {text: "cmpv.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[56:59] /* {bcV, bcS, bcK} */},
	opcmpvf64imm:              {text: "cmpv.f64@imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[16:19] /* {bcV, bcImmF64, bcK} */},
	opcmpltstr:                {text: "cmplt.str", out: bcargs[4:5] /* {bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opcmplestr:                {text: "cmple.str", out: bcargs[4:5] /* {bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opcmpgtstr:                {text: "cmpgt.str", out: bcargs[4:5] /* {bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opcmpgestr:                {text: "cmpge.str", out: bcargs[4:5] /* {bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opcmpltk:                  {text: "cmplt.k", out: bcargs[4:5] /* {bcK} */, in: bcargs[39:42] /* {bcK, bcK, bcK} */},
	opcmpltkimm:               {text: "cmplt.k@imm", out: bcargs[4:5] /* {bcK} */, in: bcargs[41:44] /* {bcK, bcImmU16, bcK} */},
	opcmplek:                  {text: "cmple.k", out: bcargs[4:5] /* {bcK} */, in: bcargs[39:42] /* {bcK, bcK, bcK} */},
	opcmplekimm:               {text: "cmple.k@imm", out: bcargs[4:5] /* {bcK} */, in: bcargs[41:44] /* {bcK, bcImmU16, bcK} */},
	opcmpgtk:                  {text: "cmpgt.k", out: bcargs[4:5] /* {bcK} */, in: bcargs[39:42] /* {bcK, bcK, bcK} */},
	opcmpgtkimm:               {text: "cmpgt.k@imm", out: bcargs[4:5] /* {bcK} */, in: bcargs[41:44] /* {bcK, bcImmU16, bcK} */},
	opcmpgek:                  {text: "cmpge.k", out: bcargs[4:5] /* {bcK} */, in: bcargs[39:42] /* {bcK, bcK, bcK} */},
	opcmpgekimm:               {text: "cmpge.k@imm", out: bcargs[4:5] /* {bcK} */, in: bcargs[41:44] /* {bcK, bcImmU16, bcK} */},
	opcmpeqf64:                {text: "cmpeq.f64", out: bcargs[4:5] /* {bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opcmpeqf64imm:             {text: "cmpeq.f64@imm", out: bcargs[4:5] /* {bcK} */, in: bcargs[106:109] /* {bcS, bcImmF64, bcK} */},
	opcmpltf64:                {text: "cmplt.f64", out: bcargs[4:5] /* {bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opcmpltf64imm:             {text: "cmplt.f64@imm", out: bcargs[4:5] /* {bcK} */, in: bcargs[106:109] /* {bcS, bcImmF64, bcK} */},
	opcmplef64:                {text: "cmple.f64", out: bcargs[4:5] /* {bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opcmplef64imm:             {text: "cmple.f64@imm", out: bcargs[4:5] /* {bcK} */, in: bcargs[106:109] /* {bcS, bcImmF64, bcK} */},
	opcmpgtf64:                {text: "cmpgt.f64", out: bcargs[4:5] /* {bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opcmpgtf64imm:             {text: "cmpgt.f64@imm", out: bcargs[4:5] /* {bcK} */, in: bcargs[106:109] /* {bcS, bcImmF64, bcK} */},
	opcmpgef64:                {text: "cmpge.f64", out: bcargs[4:5] /* {bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opcmpgef64imm:             {text: "cmpge.f64@imm", out: bcargs[4:5] /* {bcK} */, in: bcargs[106:109] /* {bcS, bcImmF64, bcK} */},
	opcmpeqi64:                {text: "cmpeq.i64", out: bcargs[4:5] /* {bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opcmpeqi64imm:             {text: "cmpeq.i64@imm", out: bcargs[4:5] /* {bcK} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opcmplti64:                {text: "cmplt.i64", out: bcargs[4:5] /* {bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opcmplti64imm:             {text: "cmplt.i64@imm", out: bcargs[4:5] /* {bcK} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opcmplei64:                {text: "cmple.i64", out: bcargs[4:5] /* {bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opcmplei64imm:             {text: "cmple.i64@imm", out: bcargs[4:5] /* {bcK} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opcmpgti64:                {text: "cmpgt.i64", out: bcargs[4:5] /* {bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opcmpgti64imm:             {text: "cmpgt.i64@imm", out: bcargs[4:5] /* {bcK} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opcmpgei64:                {text: "cmpge.i64", out: bcargs[4:5] /* {bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opcmpgei64imm:             {text: "cmpge.i64@imm", out: bcargs[4:5] /* {bcK} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opisnanf:                  {text: "isnan.f", out: bcargs[4:5] /* {bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opchecktag:                {text: "checktag", out: bcargs[10:12] /* {bcV, bcK} */, in: bcargs[53:56] /* {bcV, bcImmU16, bcK} */},
	optypebits:                {text: "typebits", out: bcargs[0:1] /* {bcS} */, in: bcargs[10:12] /* {bcV, bcK} */},
	opisnullv:                 {text: "isnull.v", out: bcargs[4:5] /* {bcK} */, in: bcargs[10:12] /* {bcV, bcK} */},
	opisnotnullv:              {text: "isnotnull.v", out: bcargs[4:5] /* {bcK} */, in: bcargs[10:12] /* {bcV, bcK} */},
	opistruev:                 {text: "istrue.v", out: bcargs[4:5] /* {bcK} */, in: bcargs[10:12] /* {bcV, bcK} */},
	opisfalsev:                {text: "isfalse.v", out: bcargs[4:5] /* {bcK} */, in: bcargs[10:12] /* {bcV, bcK} */},
	opcmpeqslice:              {text: "cmpeq.slice", out: bcargs[4:5] /* {bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opcmpeqv:                  {text: "cmpeq.v", out: bcargs[4:5] /* {bcK} */, in: bcargs[99:102] /* {bcV, bcV, bcK} */},
	opcmpeqvimm:               {text: "cmpeq.v@imm", out: bcargs[4:5] /* {bcK} */, in: bcargs[19:22] /* {bcV, bcLitRef, bcK} */},
	opdateaddmonth:            {text: "dateaddmonth", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opdateaddmonthimm:         {text: "dateaddmonth.imm", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
	opdateaddyear:             {text: "dateaddyear", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opdateaddquarter:          {text: "dateaddquarter", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opdatediffmicrosecond:     {text: "datediffmicrosecond", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opdatediffparam:           {text: "datediffparam", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[46:50] /* {bcS, bcS, bcImmU64, bcK} */},
	opdatediffmqy:             {text: "datediffmqy", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[88:92] /* {bcS, bcS, bcImmU16, bcK} */},
	opdateextractmicrosecond:  {text: "dateextractmicrosecond", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdateextractmillisecond:  {text: "dateextractmillisecond", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdateextractsecond:       {text: "dateextractsecond", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdateextractminute:       {text: "dateextractminute", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdateextracthour:         {text: "dateextracthour", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdateextractday:          {text: "dateextractday", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdateextractdow:          {text: "dateextractdow", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdateextractdoy:          {text: "dateextractdoy", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdateextractmonth:        {text: "dateextractmonth", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdateextractquarter:      {text: "dateextractquarter", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdateextractyear:         {text: "dateextractyear", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdatetounixepoch:         {text: "datetounixepoch", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdatetounixmicro:         {text: "datetounixmicro", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdatetruncmillisecond:    {text: "datetruncmillisecond", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdatetruncsecond:         {text: "datetruncsecond", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdatetruncminute:         {text: "datetruncminute", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdatetrunchour:           {text: "datetrunchour", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdatetruncday:            {text: "datetruncday", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdatetruncdow:            {text: "datetruncdow", out: bcargs[0:1] /* {bcS} */, in: bcargs[32:35] /* {bcS, bcImmU16, bcK} */},
	opdatetruncmonth:          {text: "datetruncmonth", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdatetruncquarter:        {text: "datetruncquarter", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opdatetruncyear:           {text: "datetruncyear", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opunboxts:                 {text: "unboxts", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[10:12] /* {bcV, bcK} */},
	opboxts:                   {text: "boxts", out: bcargs[10:11] /* {bcV} */, in: bcargs[3:5] /* {bcS, bcK} */, scratch: 16 * 16},
	opwidthbucketf64:          {text: "widthbucket.f64", out: bcargs[0:1] /* {bcS} */, in: bcargs[0:5] /* {bcS, bcS, bcS, bcS, bcK} */},
	opwidthbucketi64:          {text: "widthbucket.i64", out: bcargs[0:1] /* {bcS} */, in: bcargs[0:5] /* {bcS, bcS, bcS, bcS, bcK} */},
	optimebucketts:            {text: "timebucket.ts", out: bcargs[0:1] /* {bcS} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opgeohash:                 {text: "geohash", out: bcargs[0:1] /* {bcS} */, in: bcargs[1:5] /* {bcS, bcS, bcS, bcK} */, scratch: 16 * 16},
	opgeohashimm:              {text: "geohashimm", out: bcargs[0:1] /* {bcS} */, in: bcargs[88:92] /* {bcS, bcS, bcImmU16, bcK} */, scratch: 16 * 16},
	opgeotilex:                {text: "geotilex", out: bcargs[0:1] /* {bcS} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opgeotiley:                {text: "geotiley", out: bcargs[0:1] /* {bcS} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opgeotilees:               {text: "geotilees", out: bcargs[0:1] /* {bcS} */, in: bcargs[1:5] /* {bcS, bcS, bcS, bcK} */, scratch: 32 * 16},
	opgeotileesimm:            {text: "geotilees.imm", out: bcargs[0:1] /* {bcS} */, in: bcargs[88:92] /* {bcS, bcS, bcImmU16, bcK} */, scratch: 32 * 16},
	opgeodistance:             {text: "geodistance", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[0:5] /* {bcS, bcS, bcS, bcS, bcK} */},
	opalloc:                   {text: "alloc", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */, scratch: PageSize},
	opconcatstr:               {text: "concatstr", out: bcargs[3:5] /* {bcS, bcK} */, va: bcargs[3:5] /* {bcS, bcK} */, scratch: PageSize},
	opfindsym:                 {text: "findsym", out: bcargs[10:12] /* {bcV, bcK} */, in: bcargs[66:69] /* {bcB, bcSymbolID, bcK} */},
	opfindsym2:                {text: "findsym2", out: bcargs[10:12] /* {bcV, bcK} */, in: bcargs[83:88] /* {bcB, bcV, bcK, bcSymbolID, bcK} */},
	opblendv:                  {text: "blend.v", out: bcargs[10:12] /* {bcV, bcK} */, in: bcargs[70:74] /* {bcV, bcK, bcV, bcK} */},
	opblendf64:                {text: "blend.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[61:65] /* {bcS, bcK, bcS, bcK} */},
	opunpack:                  {text: "unpack", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[53:56] /* {bcV, bcImmU16, bcK} */},
	opunsymbolize:             {text: "unsymbolize", out: bcargs[10:11] /* {bcV} */, in: bcargs[10:12] /* {bcV, bcK} */},
	opunboxktoi64:             {text: "unbox.k@i64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[10:12] /* {bcV, bcK} */},
	opunboxcoercef64:          {text: "unbox.coerce.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[10:12] /* {bcV, bcK} */},
	opunboxcoercei64:          {text: "unbox.coerce.i64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[10:12] /* {bcV, bcK} */},
	opunboxcvtf64:             {text: "unbox.cvt.f64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[10:12] /* {bcV, bcK} */},
	opunboxcvti64:             {text: "unbox.cvt.i64", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[10:12] /* {bcV, bcK} */},
	opboxf64:                  {text: "box.f64", out: bcargs[10:11] /* {bcV} */, in: bcargs[3:5] /* {bcS, bcK} */, scratch: 9 * 16},
	opboxi64:                  {text: "box.i64", out: bcargs[10:11] /* {bcV} */, in: bcargs[3:5] /* {bcS, bcK} */, scratch: 9 * 16},
	opboxk:                    {text: "box.k", out: bcargs[10:11] /* {bcV} */, in: bcargs[7:9] /* {bcK, bcK} */, scratch: 16},
	opboxstr:                  {text: "box.str", out: bcargs[10:11] /* {bcV} */, in: bcargs[3:5] /* {bcS, bcK} */, scratch: PageSize},
	opboxlist:                 {text: "box.list", out: bcargs[10:11] /* {bcV} */, in: bcargs[3:5] /* {bcS, bcK} */, scratch: PageSize},
	opmakelist:                {text: "makelist", out: bcargs[10:12] /* {bcV, bcK} */, in: bcargs[4:5] /* {bcK} */, va: bcargs[10:12] /* {bcV, bcK} */, scratch: PageSize},
	opmakestruct:              {text: "makestruct", out: bcargs[10:12] /* {bcV, bcK} */, in: bcargs[4:5] /* {bcK} */, va: bcargs[96:99] /* {bcSymbolID, bcV, bcK} */, scratch: PageSize},
	ophashvalue:               {text: "hashvalue", out: bcargs[9:10] /* {bcH} */, in: bcargs[10:12] /* {bcV, bcK} */},
	ophashvalueplus:           {text: "hashvalue+", out: bcargs[9:10] /* {bcH} */, in: bcargs[9:12] /* {bcH, bcV, bcK} */},
	ophashmember:              {text: "hashmember", out: bcargs[4:5] /* {bcK} */, in: bcargs[27:30] /* {bcH, bcImmU16, bcK} */},
	ophashlookup:              {text: "hashlookup", out: bcargs[10:12] /* {bcV, bcK} */, in: bcargs[27:30] /* {bcH, bcImmU16, bcK} */},
	opaggandk:                 {text: "aggand.k", in: bcargs[38:41] /* {bcAggSlot, bcK, bcK} */},
	opaggork:                  {text: "aggor.k", in: bcargs[38:41] /* {bcAggSlot, bcK, bcK} */},
	opaggslotsumf:             {text: "aggslotsum.f64", in: bcargs[92:96] /* {bcAggSlot, bcL, bcS, bcK} */},
	opaggsumf:                 {text: "aggsum.f64", in: bcargs[50:53] /* {bcAggSlot, bcS, bcK} */},
	opaggsumi:                 {text: "aggsum.i64", in: bcargs[50:53] /* {bcAggSlot, bcS, bcK} */},
	opaggminf:                 {text: "aggmin.f64", in: bcargs[50:53] /* {bcAggSlot, bcS, bcK} */},
	opaggmini:                 {text: "aggmin.i64", in: bcargs[50:53] /* {bcAggSlot, bcS, bcK} */},
	opaggmaxf:                 {text: "aggmax.f64", in: bcargs[50:53] /* {bcAggSlot, bcS, bcK} */},
	opaggmaxi:                 {text: "aggmax.i64", in: bcargs[50:53] /* {bcAggSlot, bcS, bcK} */},
	opaggandi:                 {text: "aggand.i64", in: bcargs[50:53] /* {bcAggSlot, bcS, bcK} */},
	opaggori:                  {text: "aggor.i64", in: bcargs[50:53] /* {bcAggSlot, bcS, bcK} */},
	opaggxori:                 {text: "aggxor.i64", in: bcargs[50:53] /* {bcAggSlot, bcS, bcK} */},
	opaggcount:                {text: "aggcount", in: bcargs[38:40] /* {bcAggSlot, bcK} */},
	opaggbucket:               {text: "aggbucket", out: bcargs[6:7] /* {bcL} */, in: bcargs[36:38] /* {bcH, bcK} */},
	opaggslotandk:             {text: "aggslotand.k", in: bcargs[5:9] /* {bcAggSlot, bcL, bcK, bcK} */},
	opaggslotork:              {text: "aggslotor.k", in: bcargs[5:9] /* {bcAggSlot, bcL, bcK, bcK} */},
	opaggslotsumi:             {text: "aggslotsum.i64", in: bcargs[92:96] /* {bcAggSlot, bcL, bcS, bcK} */},
	opaggslotavgf:             {text: "aggslotavg.f64", in: bcargs[92:96] /* {bcAggSlot, bcL, bcS, bcK} */},
	opaggslotavgi:             {text: "aggslotavg.i64", in: bcargs[92:96] /* {bcAggSlot, bcL, bcS, bcK} */},
	opaggslotminf:             {text: "aggslotmin.f64", in: bcargs[92:96] /* {bcAggSlot, bcL, bcS, bcK} */},
	opaggslotmini:             {text: "aggslotmin.i64", in: bcargs[92:96] /* {bcAggSlot, bcL, bcS, bcK} */},
	opaggslotmaxf:             {text: "aggslotmax.f64", in: bcargs[92:96] /* {bcAggSlot, bcL, bcS, bcK} */},
	opaggslotmaxi:             {text: "aggslotmax.i64", in: bcargs[92:96] /* {bcAggSlot, bcL, bcS, bcK} */},
	opaggslotandi:             {text: "aggslotand.i64", in: bcargs[92:96] /* {bcAggSlot, bcL, bcS, bcK} */},
	opaggslotori:              {text: "aggslotor.i64", in: bcargs[92:96] /* {bcAggSlot, bcL, bcS, bcK} */},
	opaggslotxori:             {text: "aggslotxor.i64", in: bcargs[92:96] /* {bcAggSlot, bcL, bcS, bcK} */},
	opaggslotcount:            {text: "aggslotcount", in: bcargs[5:8] /* {bcAggSlot, bcL, bcK} */},
	opaggslotcount_v2:         {text: "aggslotcount", in: bcargs[5:8] /* {bcAggSlot, bcL, bcK} */},
	oplitref:                  {text: "litref", out: bcargs[10:11] /* {bcV} */, in: bcargs[20:21] /* {bcLitRef} */},
	opauxval:                  {text: "auxval", out: bcargs[10:12] /* {bcV, bcK} */, in: bcargs[65:66] /* {bcAuxSlot} */},
	opsplit:                   {text: "split", out: bcargs[56:59] /* {bcV, bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	optuple:                   {text: "tuple", out: bcargs[44:46] /* {bcB, bcK} */, in: bcargs[10:12] /* {bcV, bcK} */},
	opmovk:                    {text: "mov.k", out: bcargs[4:5] /* {bcK} */, in: bcargs[4:5] /* {bcK} */},
	opzerov:                   {text: "zero.v", out: bcargs[10:11] /* {bcV} */},
	opmovv:                    {text: "mov.v", out: bcargs[10:11] /* {bcV} */, in: bcargs[10:12] /* {bcV, bcK} */},
	opmovvk:                   {text: "mov.v.k", out: bcargs[10:12] /* {bcV, bcK} */, in: bcargs[10:12] /* {bcV, bcK} */},
	opmovf64:                  {text: "mov.f64", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opmovi64:                  {text: "mov.i64", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opobjectsize:              {text: "objectsize", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[10:12] /* {bcV, bcK} */},
	oparraysize:               {text: "arraysize", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	oparrayposition:           {text: "arrayposition", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[69:72] /* {bcS, bcV, bcK} */},
	opCmpStrEqCs:              {text: "cmp_str_eq_cs", out: bcargs[4:5] /* {bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opCmpStrEqCi:              {text: "cmp_str_eq_ci", out: bcargs[4:5] /* {bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opCmpStrEqUTF8Ci:          {text: "cmp_str_eq_utf8_ci", out: bcargs[4:5] /* {bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opCmpStrFuzzyA3:           {text: "cmp_str_fuzzy_A3", out: bcargs[4:5] /* {bcK} */, in: bcargs[22:26] /* {bcS, bcS, bcDictSlot, bcK} */},
	opCmpStrFuzzyUnicodeA3:    {text: "cmp_str_fuzzy_unicode_A3", out: bcargs[4:5] /* {bcK} */, in: bcargs[22:26] /* {bcS, bcS, bcDictSlot, bcK} */},
	opHasSubstrFuzzyA3:        {text: "contains_fuzzy_A3", out: bcargs[4:5] /* {bcK} */, in: bcargs[22:26] /* {bcS, bcS, bcDictSlot, bcK} */},
	opHasSubstrFuzzyUnicodeA3: {text: "contains_fuzzy_unicode_A3", out: bcargs[4:5] /* {bcK} */, in: bcargs[22:26] /* {bcS, bcS, bcDictSlot, bcK} */},
	opSkip1charLeft:           {text: "skip_1char_left", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opSkip1charRight:          {text: "skip_1char_right", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opSkipNcharLeft:           {text: "skip_nchar_left", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opSkipNcharRight:          {text: "skip_nchar_right", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[2:5] /* {bcS, bcS, bcK} */},
	opTrimWsLeft:              {text: "trim_ws_left", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opTrimWsRight:             {text: "trim_ws_right", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opTrim4charLeft:           {text: "trim_char_left", out: bcargs[0:1] /* {bcS} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opTrim4charRight:          {text: "trim_char_right", out: bcargs[0:1] /* {bcS} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opoctetlength:             {text: "octetlength", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opcharlength:              {text: "characterlength", out: bcargs[0:1] /* {bcS} */, in: bcargs[3:5] /* {bcS, bcK} */},
	opSubstr:                  {text: "substr", out: bcargs[0:1] /* {bcS} */, in: bcargs[1:5] /* {bcS, bcS, bcS, bcK} */},
	opSplitPart:               {text: "split_part", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[59:63] /* {bcS, bcDictSlot, bcS, bcK} */},
	opContainsPrefixCs:        {text: "contains_prefix_cs", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opContainsPrefixCi:        {text: "contains_prefix_ci", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opContainsPrefixUTF8Ci:    {text: "contains_prefix_utf8_ci", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opContainsSuffixCs:        {text: "contains_suffix_cs", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opContainsSuffixCi:        {text: "contains_suffix_ci", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opContainsSuffixUTF8Ci:    {text: "contains_suffix_utf8_ci", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opContainsSubstrCs:        {text: "contains_substr_cs", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opContainsSubstrCi:        {text: "contains_substr_ci", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opContainsSubstrUTF8Ci:    {text: "contains_substr_utf8_ci", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opEqPatternCs:             {text: "eq_pattern_cs", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opEqPatternCi:             {text: "eq_pattern_ci", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opEqPatternUTF8Ci:         {text: "eq_pattern_utf8_ci", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opContainsPatternCs:       {text: "contains_pattern_cs", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opContainsPatternCi:       {text: "contains_pattern_ci", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opContainsPatternUTF8Ci:   {text: "contains_pattern_utf8_ci", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opIsSubnetOfIP4:           {text: "is_subnet_of_ip4", out: bcargs[4:5] /* {bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opDfaT6:                   {text: "dfa_tiny6", out: bcargs[4:5] /* {bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opDfaT7:                   {text: "dfa_tiny7", out: bcargs[4:5] /* {bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opDfaT8:                   {text: "dfa_tiny8", out: bcargs[4:5] /* {bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opDfaT6Z:                  {text: "dfa_tiny6Z", out: bcargs[4:5] /* {bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opDfaT7Z:                  {text: "dfa_tiny7Z", out: bcargs[4:5] /* {bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opDfaT8Z:                  {text: "dfa_tiny8Z", out: bcargs[4:5] /* {bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opDfaLZ:                   {text: "dfa_largeZ", out: bcargs[4:5] /* {bcK} */, in: bcargs[23:26] /* {bcS, bcDictSlot, bcK} */},
	opslower:                  {text: "slower", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */, scratch: PageSize},
	opsupper:                  {text: "supper", out: bcargs[3:5] /* {bcS, bcK} */, in: bcargs[3:5] /* {bcS, bcK} */, scratch: PageSize},
	opaggapproxcount:          {text: "aggapproxcount", in: bcargs[26:30] /* {bcAggSlot, bcH, bcImmU16, bcK} */},
	opaggapproxcountmerge:     {text: "aggapproxcountmerge", in: bcargs[102:106] /* {bcAggSlot, bcS, bcImmU16, bcK} */},
	opaggslotapproxcount:      {text: "aggslotapproxcount", in: bcargs[78:83] /* {bcAggSlot, bcL, bcH, bcImmU16, bcK} */},
	opaggslotapproxcountmerge: {text: "aggslotapproxcountmerge", in: bcargs[30:35] /* {bcAggSlot, bcL, bcS, bcImmU16, bcK} */},
	oppowuintf64:              {text: "powuint.f64", out: bcargs[0:1] /* {bcS} */, in: bcargs[13:16] /* {bcS, bcImmI64, bcK} */},
}

var bcargs = [109]bcArgType{bcS, bcS, bcS, bcS, bcK, bcAggSlot, bcL, bcK,
	bcK, bcH, bcV, bcK, bcS, bcS, bcImmI64, bcK, bcV, bcImmF64, bcK,
	bcV, bcLitRef, bcK, bcS, bcS, bcDictSlot, bcK, bcAggSlot, bcH,
	bcImmU16, bcK, bcAggSlot, bcL, bcS, bcImmU16, bcK, bcB, bcH, bcK,
	bcAggSlot, bcK, bcK, bcK, bcImmU16, bcK, bcB, bcK, bcS, bcS,
	bcImmU64, bcK, bcAggSlot, bcS, bcK, bcV, bcImmU16, bcK, bcV, bcS,
	bcK, bcS, bcDictSlot, bcS, bcK, bcS, bcK, bcAuxSlot, bcB,
	bcSymbolID, bcK, bcS, bcV, bcK, bcV, bcK, bcK, bcV, bcImmI64, bcK,
	bcAggSlot, bcL, bcH, bcImmU16, bcK, bcB, bcV, bcK, bcSymbolID, bcK,
	bcS, bcS, bcImmU16, bcK, bcAggSlot, bcL, bcS, bcK, bcSymbolID, bcV,
	bcK, bcV, bcV, bcK, bcAggSlot, bcS, bcImmU16, bcK, bcS, bcImmF64,
	bcK}
