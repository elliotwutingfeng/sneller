// Copyright (C) 2023 Sneller, Inc.
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

TEXT asinf32(SB), NOSPLIT|NOFRAME, $0
  //IN Z0 = 16x float32
  //IN K1 = active mask
  VEXTRACTF32X8  $1, Z0, Y3
  VCVTPS2PD Y0, Z2
  VCVTPS2PD Y3, Z3
  KSHIFTRW $8, K1, K2

  KTESTW K1, K1
  JZ next

  VBROADCASTSD CONSTF64_ABS_BITS(), Z7
  VBROADCASTSD CONSTF64_HALF(), Z8
  VBROADCASTSD CONSTF64_1(), Z9

  VANDPD.Z Z7, Z2, K1, Z6
  VANDPD.Z Z7, Z3, K2, Z7
  VCMPPD $VCMP_IMM_LT_OQ, Z8, Z6, K1, K3
  VCMPPD $VCMP_IMM_LT_OQ, Z8, Z7, K2, K4
  VSUBPD Z6, Z9, Z4
  VSUBPD Z7, Z9, Z5
  VMULPD Z8, Z4, Z4
  VMULPD Z8, Z5, Z5
  VMULPD Z2, Z2, K3, Z4
  VMULPD Z3, Z3, K4, Z5
  VSQRTPD Z4, Z12
  VSQRTPD Z5, Z13
  VMULPD Z12, Z12, Z14
  VMULPD Z13, Z13, Z15
  VMOVAPD Z12, Z16
  VMOVAPD Z13, Z17
  VFMSUB213PD Z14, Z12, Z16 // Z16 = (Z12 * Z16) - Z14
  VFMSUB213PD Z15, Z13, Z17 // Z17 = (Z13 * Z17) - Z15
  VADDPD Z14, Z4, Z18
  VADDPD Z15, Z5, Z19
  VSUBPD Z4, Z18, Z20
  VSUBPD Z5, Z19, Z21
  VSUBPD Z20, Z18, Z22
  VSUBPD Z21, Z19, Z23
  VSUBPD Z22, Z4, Z22
  VSUBPD Z23, Z5, Z23
  VSUBPD Z20, Z14, Z14
  VSUBPD Z21, Z15, Z15
  VADDPD Z22, Z14, Z14
  VADDPD Z23, Z15, Z15
  VADDPD Z14, Z16, Z14
  VADDPD Z15, Z17, Z15
  VDIVPD Z12, Z9, Z16
  VDIVPD Z13, Z9, Z17
  VFNMADD213PD Z9, Z16, Z12 // Z12 = -(Z16 * Z12) + Z9
  VFNMADD213PD Z9, Z17, Z13 // Z13 = -(Z17 * Z13) + Z9
  VMULPD Z12, Z16, Z12
  VMULPD Z13, Z17, Z13
  VMULPD Z18, Z16, Z20
  VMULPD Z19, Z17, Z21
  VMOVAPD Z16, Z22
  VMOVAPD Z17, Z23
  VFMSUB213PD Z20, Z18, Z22 // Z22 = (Z18 * Z22) - Z20
  VFMSUB213PD Z21, Z19, Z23 // Z23 = (Z19 * Z23) - Z21
  VFMADD231PD Z14, Z16, Z22 // Z22 = (Z16 * Z14) + Z22
  VFMADD231PD Z15, Z17, Z23 // Z23 = (Z17 * Z15) + Z23
  VFMADD231PD Z12, Z18, Z22 // Z22 = (Z18 * Z12) + Z22
  VFMADD231PD Z13, Z19, Z23 // Z23 = (Z19 * Z13) + Z23
  VMULPD Z8, Z20, Z12
  VMULPD Z8, Z21, Z13
  VMOVAPD Z6, K3, Z12
  VMOVAPD Z7, K4, Z13
  VCMPPD $VCMP_IMM_EQ_OQ, Z9, Z6, K5
  VCMPPD $VCMP_IMM_EQ_OQ, Z9, Z7, K6
  VXORPD Z12, Z12, K5, Z12
  VXORPD Z13, Z13, K6, Z13
  KORW K3, K5, K5
  KORW K4, K6, K6
  KNOTW K5, K5
  KNOTW K6, K6
  VMULPD.Z Z8, Z22, K5, Z6
  VMULPD.Z Z8, Z23, K6, Z7
  VMULPD Z4, Z4, Z8
  VMULPD Z5, Z5, Z9
  VMULPD Z8, Z8, Z10
  VMULPD Z9, Z9, Z11
  VMULPD Z10, Z10, Z14
  VMULPD Z11, Z11, Z15
  VBROADCASTSD CONST_GET_PTR(const_asin, 0), Z16
  VBROADCASTSD CONST_GET_PTR(const_asin, 16), Z18
  VBROADCASTSD CONST_GET_PTR(const_asin, 8), Z20
  VBROADCASTSD CONST_GET_PTR(const_asin, 24), Z21
  VMOVAPD Z16, Z17
  VFMADD213PD Z20, Z4, Z16 // Z16 = (Z4 * Z16) + Z20
  VFMADD213PD Z20, Z5, Z17 // Z17 = (Z5 * Z17) + Z20
  VMOVAPD Z18, Z19
  VFMADD213PD Z21, Z4, Z18 // Z18 = (Z4 * Z18) + Z21
  VFMADD213PD Z21, Z5, Z19 // Z19 = (Z5 * Z19) + Z21
  VBROADCASTSD CONST_GET_PTR(const_asin, 32), Z20
  VBROADCASTSD CONST_GET_PTR(const_asin, 48), Z22
  VMOVAPD Z20, Z21
  VMOVAPD Z22, Z23
  VBROADCASTSD CONST_GET_PTR(const_asin, 40), Z24
  VBROADCASTSD CONST_GET_PTR(const_asin, 56), Z25
  VFMADD213PD Z24, Z4, Z20 // Z20 = (Z4 * Z20) + Z24
  VFMADD213PD Z24, Z5, Z21 // Z21 = (Z5 * Z21) + Z24
  VFMADD213PD Z25, Z4, Z22 // Z22 = (Z4 * Z22) + Z25
  VFMADD213PD Z25, Z5, Z23 // Z23 = (Z5 * Z23) + Z25
  VFMADD231PD Z16, Z8, Z18 // Z18 = (Z8 * Z16) + Z18
  VFMADD231PD Z17, Z9, Z19 // Z19 = (Z9 * Z17) + Z19
  VFMADD231PD Z20, Z8, Z22 // Z22 = (Z8 * Z20) + Z22
  VFMADD231PD Z21, Z9, Z23 // Z23 = (Z9 * Z21) + Z23
  VBROADCASTSD CONST_GET_PTR(const_asin, 64), Z16
  VBROADCASTSD CONST_GET_PTR(const_asin, 80), Z20
  VMOVAPD Z16, Z17
  VMOVAPD Z20, Z21
  VBROADCASTSD CONST_GET_PTR(const_asin, 72), Z24
  VBROADCASTSD CONST_GET_PTR(const_asin, 88), Z25
  VFMADD213PD Z24, Z4, Z16 // Z16 = (Z4 * Z16) + Z24
  VFMADD213PD Z24, Z5, Z17 // Z17 = (Z5 * Z17) + Z24
  VFMADD213PD Z25, Z4, Z20 // Z20 = (Z4 * Z20) + Z25
  VFMADD213PD Z25, Z5, Z21 // Z21 = (Z5 * Z21) + Z25
  VFMADD231PD Z16, Z8, Z20  // Z20 = (Z8 * Z16) + Z20
  VFMADD231PD Z17, Z9, Z21  // Z21 = (Z9 * Z17) + Z21
  VFMADD231PD Z22, Z10, Z20 // Z20 = (Z10 * Z22) + Z20
  VFMADD231PD Z23, Z11, Z21 // Z21 = (Z11 * Z23) + Z21
  VFMADD231PD Z18, Z14, Z20 // Z20 = (Z14 * Z18) + Z20
  VFMADD231PD Z19, Z15, Z21 // Z21 = (Z15 * Z19) + Z21
  VMULPD Z12, Z4, Z4
  VMULPD Z13, Z5, Z5
  VBROADCASTSD CONST_GET_PTR(const_asin, 96), Z9
  VBROADCASTSD CONST_GET_PTR(const_asin, 104), Z14
  VMULPD Z4, Z20, Z4
  VMULPD Z5, Z21, Z5
  VSUBPD Z12, Z9, Z10
  VSUBPD Z13, Z9, Z11
  VSUBPD Z10, Z9, Z8
  VSUBPD Z11, Z9, Z9
  VSUBPD Z12, Z8, Z8
  VSUBPD Z13, Z9, Z9
  VADDPD Z14, Z8, Z8
  VADDPD Z14, Z9, Z9
  VSUBPD Z6, Z8, Z6
  VSUBPD Z7, Z9, Z7
  VSUBPD Z4, Z10, Z8
  VSUBPD Z5, Z11, Z9
  VSUBPD Z8, Z10, Z10
  VSUBPD Z9, Z11, Z11
  VSUBPD Z4, Z10, Z10
  VSUBPD Z5, Z11, Z11
  VADDPD Z6, Z10, Z6
  VADDPD Z7, Z11, Z7
  VADDPD Z6, Z8, Z6
  VADDPD Z7, Z9, Z7
  VADDPD Z6, Z6, Z6
  VADDPD Z7, Z7, Z7
  VBROADCASTSD CONSTF64_SIGN_BIT(), Z10
  VADDPD Z4, Z12, K3, Z6
  VADDPD Z5, Z13, K4, Z7
  VPTERNLOGQ $108, Z10, Z6, K1, Z2
  VPTERNLOGQ $108, Z10, Z7, K2, Z3

next:

//OUT Z2 = 16 float32
//OUT K1 = active lanes
    VCVTPD2PS Z2, Y0
	VCVTPD2PS Z3, Y3
	VINSERTF32X8 $1, Y3, Z0, Z0
    RET
//  BC_UNPACK_2xSLOT(0, OUT(DX), OUT(R8))
//  BC_STORE_F64_TO_SLOT(IN(Z2), IN(Z3), IN(DX))
//  BC_STORE_K_TO_SLOT(IN(K1), IN(R8))
//  NEXT_ADVANCE(BC_SLOT_SIZE*4)
