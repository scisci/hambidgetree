package hambidgetree

import "math"

var GoldenRatio = math.Sqrt(5)/2 + 0.5 // Equivalent to math.Phi
var SquareRatio = 1.0
var RootFiveRatio = math.Sqrt(5)

/*
func NewGoldenRatios() Ratios {
	return NewRatios([]float64{
		0.125, // 0.125 1/2 of 0.25
		1 / (math.Sqrt(5) + 5),              // 0.138 1/2 of 0.276
		2 / (math.Sqrt(5)*3 + 7),            // 0.146 1/2 of 0.2918
		1 / (2*math.Sqrt(5) + 2),            // 0.154 1/2 of 0.309
		math.Sqrt(5) / (4 + 4*math.Sqrt(5)), //0.172 1/2 of 0.3455
		2 / (9 + math.Sqrt(5)),              // 0.178 1/2 of 0.3559
		1 / (10 - 2*math.Sqrt(5)),           //0.181 1/2 of 0.3618
		1 / (math.Sqrt(5) + 3),              // 0.191 1/2 of 0.382
		1 / (math.Sqrt(5)*4 - 4),            // 0.202 1/2 of 0.4045
		math.Sqrt(5) / (2 + 4*math.Sqrt(5)), //0.204 1/2 of 0.408
		2 / (math.Sqrt(5) + 7),              // 0.216 1/2 of 0.433
		1 / (2 * math.Sqrt(5)),              // 0.2236 1/2 of 0.4472
		0.25, // 0.25
		1 / (math.Sqrt(5)*3 - 3), // 0.2696 1/2 of 0.5393

		2 / (math.Sqrt(5) + 5),              // 0.2764
		4 / (math.Sqrt(5)*3 + 7),            // 0.2918
		1 / (math.Sqrt(5) + 1),              // 0.309
		math.Sqrt(5) / (2 + 2*math.Sqrt(5)), // 0.3455
		4 / (9 + math.Sqrt(5)),              // 0.3559
		1 / (5 - math.Sqrt(5)),              // 0.3618
		2 / (math.Sqrt(5) + 3),              // 0.382
		1 / (math.Sqrt(5)*2 - 2),            // 0.4045
		math.Sqrt(5) / (1 + 2*math.Sqrt(5)), // 0.408
		//2 / (7 - math.Sqrt(5)),    // 0.4198 1/2 of 0.8396
		4 / (math.Sqrt(5) + 7), // 0.433
		1 / math.Sqrt(5),       // 0.4472
		0.5,                    // 0.5
		2 / (math.Sqrt(5)*3 - 3), // 0.5393
		4 / (math.Sqrt(5) + 5),   // 0.5528
		//math.Sqrt(5) / 4,          // 0.559 1/2 of 1.118
		//8 / (math.Sqrt(5) * 3 + 7), //0.5835 inverse of 1.17135
		0.875 - math.Sqrt(5)/8, // 0.5955 1/2 of 1.191
		2 / (math.Sqrt(5) + 1), // 0.618
		//math.Sqrt(5) / 8 + 0.375,  // 0.6545 1/2 of 1.309
		math.Sqrt(5) / (1 + math.Sqrt(5)), // 0.691
		//8 / (math.Sqrt(5) + 9),  // 0.712 inverse of 1.4045
		2 / (5 - math.Sqrt(5)), // 0.7236
		4 / (math.Sqrt(5) + 3), // 0.764
		1 / (math.Sqrt(5) - 1), // 0.809
		//(math.Sqrt(5) * 2) / (math.Sqrt(5) * 2 + 1),// 0.817 inverse of 1.2236
		4 / (7 - math.Sqrt(5)), // 0.8396
		8 / (math.Sqrt(5) + 7), // 0.866 inverse of 1.1545
		2 / math.Sqrt(5),       // 0.894
		//math.Sqrt(5) / 8 + 0.625,  // 0.9045 1/2 of 1.809
		//0.75 * (math.Sqrt(5) - 1), // 0.927 1/2 of 1.854
		1, // 1
		//4 / (math.Sqrt(5) * 3 - 3), // 1.078 inverse of .927
		//8 / (math.Sqrt(5) + 5),  // 1.105 inverse of 0.9045
		math.Sqrt(5) / 2,       // 1.118
		math.Sqrt(5)/8 + 0.875, // 1.1545
		1.75 - math.Sqrt(5)/4,  // 1.191,
		//1 / (math.Sqrt(5) * 2) + 1,// 1.2236 1/2 of 2.4472
		math.Sqrt(5) - 1,      // 1.236
		math.Sqrt(5)/4 + 0.75, // 1.309
		2.5 - math.Sqrt(5)/2,  // 1.382,
		//math.Sqrt(5) / 8 + 1.125,  // 1.4045 1/2 of 2.809
		1/math.Sqrt(5) + 1, // 1.4472
		//8 / (math.Sqrt(5) + 3),  // 1.527 inverse of 0.6545
		math.Sqrt(5)/2 + 0.5,   // 1.618
		8 / (7 - math.Sqrt(5)), // 1.679 inverse of 0.5955
		//math.Sqrt(5) * 0.375 + 0.875,//1.7135 1/2 of 3.427
		//4 / math.Sqrt(5),          // 1.788 inverse of 0.559
		math.Sqrt(5)/4 + 1.25,    // 1.809
		1.5 * (math.Sqrt(5) - 1), // 1.854
		2,                     // 2
		math.Sqrt(5),          // 2.236
		math.Sqrt(5)/4 + 1.75, // 2.309
		//(7 - math.Sqrt(5)) / 2,    // 2.382 inverse of 0.4198
		1/math.Sqrt(5) + 2,       // 2.4472
		math.Sqrt(5)*2 - 2,       // 2.472
		math.Sqrt(5)/2 + 1.5,     // 2.618
		5 - math.Sqrt(5),         // 2.764
		math.Sqrt(5)/4 + 2.25,    // 2.809
		2/math.Sqrt(5) + 2,       // 2.8944
		math.Sqrt(5) + 1,         // 3.236
		math.Sqrt(5)*0.75 + 1.75, // 3.427
		math.Sqrt(5)/2 + 2.5,     // 3.618

		(math.Sqrt(5)*3 - 3),                // 3.709 inverse of 0.2696
		4,                                   // 4
		(2 * math.Sqrt(5)),                  // 4.472 inverse of .2236
		(math.Sqrt(5) + 7) / 2,              // 4.618 inverse of .216
		(2 + 4*math.Sqrt(5)) / math.Sqrt(5), // 4.894 inverse of 0.204
		(math.Sqrt(5)*4 - 4),                // 4.944 inverse of 0.202
		(math.Sqrt(5) + 3),                  // 5.23 inverse of 0.191
		(10 - 2*math.Sqrt(5)),               //5.527 inverse of 0.181
		(9 + math.Sqrt(5)) / 2,              // 5.618 inverse of 0.178
		(4 + 4*math.Sqrt(5)) / math.Sqrt(5), //5.788 inverse of 0.172
		(2*math.Sqrt(5) + 2),                // 6.47 inverse of 0.154
		(math.Sqrt(5)*3 + 7) / 2,            // 6.854 inverse of 0.146
		(math.Sqrt(5) + 5),                  // 7.23 inverse of 0.138
		8,
	})
}
*/

func NewGoldenRatios() Ratios {
	return NewExprRatios([]string{
		"1/8",                   // 0.125 1/2 of 0.25
		"1/(SQRT(5)+5)",         // 0.138 1/2 of 0.276
		"2/(SQRT(5)*3+7)",       // 0.146 1/2 of 0.2918
		"1/(2*SQRT(5)+2)",       // 0.154 1/2 of 0.309
		"SQRT(5)/(4+4*SQRT(5))", //0.172 1/2 of 0.3455
		"2/(9+SQRT(5))",         // 0.178 1/2 of 0.3559
		"1/(10-2*SQRT(5))",      //0.181 1/2 of 0.3618
		"1/(SQRT(5)+3)",         // 0.191 1/2 of 0.382
		"1/(SQRT(5)*4-4)",       // 0.202 1/2 of 0.4045
		"SQRT(5)/(2+4*SQRT(5))", //0.204 1/2 of 0.408
		"2/(SQRT(5)+7)",         // 0.216 1/2 of 0.433
		"1/(2*SQRT(5))",         // 0.2236 1/2 of 0.4472
		"1/4",                   // 0.25
		"1/(SQRT(5)*3-3)",       // 0.2696 1/2 of 0.5393

		"2/(SQRT(5)+5)",         // 0.2764
		"4/(SQRT(5)*3+7)",       // 0.2918
		"1/(SQRT(5)+1)",         // 0.309
		"SQRT(5)/(2+2*SQRT(5))", // 0.3455
		"4/(9+SQRT(5))",         // 0.3559
		"1/(5-SQRT(5))",         // 0.3618
		"2/(SQRT(5)+3)",         // 0.382
		"1/(SQRT(5)*2-2)",       // 0.4045
		"SQRT(5)/(1+2*SQRT(5))", // 0.408
		//2 / (7 - SQRT(5)),    // 0.4198 1/2 of 0.8396
		"4/(SQRT(5)+7)",   // 0.433
		"1/SQRT(5)",       // 0.4472
		"1/2",             // 0.5
		"2/(SQRT(5)*3-3)", // 0.5393
		"4/(SQRT(5)+5)",   // 0.5528
		//SQRT(5) / 4,          // 0.559 1/2 of 1.118
		//8 / (SQRT(5) * 3 + 7), //0.5835 inverse of 1.17135
		"(7-SQRT(5))/8", // 0.5955 1/2 of 1.191
		"2/(SQRT(5)+1)", // 0.618
		//SQRT(5) / 8 + 0.375,  // 0.6545 1/2 of 1.309
		"SQRT(5)/(1+SQRT(5))", // 0.691
		//8 / (SQRT(5) + 9),  // 0.712 inverse of 1.4045
		"2/(5-SQRT(5))", // 0.7236
		"4/(SQRT(5)+3)", // 0.764
		"1/(SQRT(5)-1)", // 0.809
		//(SQRT(5) * 2) / (SQRT(5) * 2 + 1),// 0.817 inverse of 1.2236
		"4/(7-SQRT(5))", // 0.8396
		"8/(SQRT(5)+7)", // 0.866 inverse of 1.1545
		"2/SQRT(5)",     // 0.894
		//SQRT(5) / 8 + 0.625,  // 0.9045 1/2 of 1.809
		//0.75 * (SQRT(5) - 1), // 0.927 1/2 of 1.854
		"1", // 1
		//4 / (SQRT(5) * 3 - 3), // 1.078 inverse of .927
		//8 / (SQRT(5) + 5),  // 1.105 inverse of 0.9045
		"SQRT(5)/2",     // 1.118
		"(SQRT(5)+7)/8", // 1.1545
		"(7-SQRT(5))/4", // 1.191,
		//1 / (SQRT(5) * 2) + 1,// 1.2236 1/2 of 2.4472
		"SQRT(5)-1",     // 1.236
		"(SQRT(5)+3)/4", // 1.309
		"(5-SQRT(5))/2", // 1.382,
		//SQRT(5) / 8 + 1.125,  // 1.4045 1/2 of 2.809
		"1/SQRT(5)+1", // 1.4472
		//8 / (SQRT(5) + 3),  // 1.527 inverse of 0.6545
		"(SQRT(5)+1)/2", // 1.618
		"8/(7-SQRT(5))", // 1.679 inverse of 0.5955
		//SQRT(5) * 0.375 + 0.875,//1.7135 1/2 of 3.427
		//4 / SQRT(5),          // 1.788 inverse of 0.559
		"(SQRT(5)+5)/4",   // 1.809
		"(3*SQRT(5)-3)/2", // 1.854
		"2",               // 2
		"SQRT(5)",         // 2.236
		"(7+SQRT(5))/4",   // 2.309
		//(7 - SQRT(5)) / 2,    // 2.382 inverse of 0.4198
		"1/SQRT(5)+2",     // 2.4472
		"SQRT(5)*2-2",     // 2.472
		"(SQRT(5)+3)/2",   // 2.618
		"5-SQRT(5)",       // 2.764
		"(SQRT(5)+9)/4",   // 2.809
		"2/SQRT(5)+2",     // 2.8944
		"SQRT(5)+1",       // 3.236
		"(3*SQRT(5)+7)/4", // 3.427
		"(5+SQRT(5))/2",   // 3.618

		"(SQRT(5)*3-3)",         // 3.709 inverse of 0.2696
		"4",                     // 4
		"(2*SQRT(5))",           // 4.472 inverse of .2236
		"(SQRT(5)+7) / 2",       // 4.618 inverse of .216
		"(2+4*SQRT(5))/SQRT(5)", // 4.894 inverse of 0.204
		"(SQRT(5)*4-4)",         // 4.944 inverse of 0.202
		"(SQRT(5)+3)",           // 5.23 inverse of 0.191
		"(10-2*SQRT(5))",        //5.527 inverse of 0.181
		"(9+SQRT(5))/2",         // 5.618 inverse of 0.178
		"(4+4*SQRT(5))/SQRT(5)", //5.788 inverse of 0.172
		"(2*SQRT(5)+2)",         // 6.47 inverse of 0.154
		"(SQRT(5)*3+7)/2",       // 6.854 inverse of 0.146
		"(SQRT(5)+5)",           // 7.23 inverse of 0.138
		"8",
	})
}
