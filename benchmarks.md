# Benchmarks:

    BenchmarkSearchValid-4                    50      25884559 ns/op
    BenchmarkSearchValidSlow-4                50      27821677 ns/op
    BenchmarkMidGamePruningDepth2-4           10     138513378 ns/op
    BenchmarkOpeningPruningDepth2-4           20      94688395 ns/op
    BenchmarkOpeningPruningDepth3-4            1	1441998152 ns/op
    BenchmarkMidGamePruningDepth3-4            1	2177541792 ns/op
    BenchmarkOpeningPruningDepth4-4            1	16566151366 ns/op
    BenchmarkMidGamePruningDepth4-4            1	16079072907 ns/op


After I change the []byte slice board to a [120]byte array, and don't copy it:

    BenchmarkSearchValid-4                       100      22877600 ns/op
    BenchmarkSearchValidSlow-4                    50      29033893 ns/op
    BenchmarkMidGamePruningDepth2-4               10     136505438 ns/op
    BenchmarkOpeningPruningDepth2-4               20      77758483 ns/op
    BenchmarkOpeningPruningDepth3-4                1	1257017288 ns/op
    BenchmarkMidGamePruningDepth3-4                1	2254520731 ns/op
    BenchmarkMidGameTwoPruningDepth3-4        300000          6268 ns/op
    BenchmarkOpeningOrderedDepth3-4                1	1341534583 ns/op
    BenchmarkMidGameOrderedDepth3-4                1	2325314282 ns/op
    BenchmarkMidGameTwoOrderedDepth3-4        200000          6107 ns/op
    BenchmarkOpeningPruningDepth4-4                1	15881901832 ns/op
    BenchmarkMidGamePruningDepth4-4                1	18561026485 ns/op
    PASS

Benchmarks after I figured out that I wasn't calling MiniMaxOrdered inside of MinimaxOrdered...

    BenchmarkMidGamePruningDepth2-4               10     136375546 ns/op
    BenchmarkOpeningPruningDepth2-4               20      91753078 ns/op
    BenchmarkOpeningPruningDepth3-4                1	1322371490 ns/op
    BenchmarkMidGamePruningDepth3-4                1	2164638763 ns/op
    BenchmarkMidGameTwoPruningDepth3-4        300000          6137 ns/op
    BenchmarkOpeningOrderedDepth3-4                1	1246987176 ns/op
    BenchmarkMidGameOrderedDepth3-4                1	2455577971 ns/op
    BenchmarkMidGameTwoOrderedDepth3-4        300000          6135 ns/op
    BenchmarkOpeningPruningDepth4-4                1	15661720638 ns/op
    BenchmarkMidGamePruningDepth4-4                1	18284754487 ns/op

Giving up with PV:

    BenchmarkMidGamePruningDepth2-4               10     131262359 ns/op
    BenchmarkOpeningPruningDepth2-4               20      94373512 ns/op
    BenchmarkOpeningPruningDepth3-4                1	1350511232 ns/op
    BenchmarkMidGamePruningDepth3-4                1	2508213115 ns/op
    BenchmarkOpeningPruningDepth4-4                1	16827821614 ns/op
    BenchmarkMidGamePruningDepth4-4                1	15570438668 ns/op
    PASS
    ok      github.com/polypmer/ghess	50.129s

After Using profiling

    BenchmarkSearchValid-4                       100      15393451 ns/op
    BenchmarkSearchValidSlow-4                   100      15275635 ns/op
    BenchmarkMidGamePruningDepth2-4               20      93899457 ns/op
    BenchmarkOpeningPruningDepth2-4               20      56920214 ns/op
    BenchmarkOpeningPruningDepth3-4                2     784527542 ns/op
    BenchmarkMidGamePruningDepth3-4                1	1518414649 ns/op
    BenchmarkMidGameTwoPruningDepth3-4             3     361195704 ns/op
    BenchmarkOpeningPruningDepth4-4                1	10112804027 ns/op
    BenchmarkMidGamePruningDepth4-4                1	9730502681 ns/op
    PASS
    ok      github.com/polypmer/ghess	31.456s


Benchmarks with pawn validation in Check reduced...

    BenchmarkSearchValid-4                       100      15480015 ns/op
    BenchmarkSearchValidSlow-4                   100      15150476 ns/op
    BenchmarkMidGamePruningDepth2-4               20      91290571 ns/op
    BenchmarkOpeningPruningDepth2-4               30      50329151 ns/op
    BenchmarkOpeningPruningDepth3-4                2     686627412 ns/op
    BenchmarkMidGamePruningDepth3-4                1	1249006486 ns/op
    BenchmarkMidGameTwoPruningDepth3-4             5     342692638 ns/op
    BenchmarkOpeningPruningDepth4-4                1	8347604470 ns/op
    BenchmarkMidGamePruningDepth4-4                1	8536214158 ns/op
    PASS
    ok      github.com/polypmer/ghess	30.632s

Benchmark After certain profiling:

    BenchmarkSearchValid-4                       100      12357042 ns/op
    BenchmarkSearchValidSlow-4                   100      14710691 ns/op
    BenchmarkMidGamePruningDepth2-4               20      78253547 ns/op
    BenchmarkOpeningPruningDepth2-4               30      44239677 ns/op
    BenchmarkOpeningPruningDepth3-4                2     617343700 ns/op
    BenchmarkMidGamePruningDepth3-4                1	1478115205 ns/op
    BenchmarkMidGameTwoPruningDepth3-4             3     342789913 ns/op
    BenchmarkOpeningPruningDepth4-4                1	9527661212 ns/op
    BenchmarkMidGamePruningDepth4-4                1	10551483526 ns/op
    PASS
    ok      github.com/polypmer/ghess	30.863s

Without Updating Check within the Move method

    BenchmarkSearchValid-4                       100      10165640 ns/op
    BenchmarkSearchValidSlow-4                   100      10950169 ns/op
    BenchmarkMidGamePruningDepth2-4               20      55903070 ns/op
    BenchmarkOpeningPruningDepth2-4               50      36001671 ns/op
    BenchmarkOpeningPruningDepth3-4                2     500474578 ns/op
    BenchmarkMidGamePruningDepth3-4                1	1110722333 ns/op
    BenchmarkMidGameTwoPruningDepth3-4             5     239901689 ns/op
    BenchmarkOpeningPruningDepth4-4                1	6221718962 ns/op
    BenchmarkMidGamePruningDepth4-4                1	6580892546 ns/op
    PASS
    ok      github.com/polypmer/ghess	22.101s

Before New Check Method:

    BenchmarkSearchValid-4                      5000        470955 ns/op
    BenchmarkSearchValidSlow-4                  1000       1380257 ns/op
    BenchmarkMidGamePruningDepth2-4               50      42497511 ns/op
    BenchmarkOpeningPruningDepth2-4              100      26641277 ns/op
    BenchmarkOpeningPruningDepth3-4               10     246942839 ns/op
    BenchmarkMidGamePruningDepth3-4                5     279948987 ns/op
    BenchmarkMidGameTwoPruningDepth3-4            20     133189782 ns/op
    BenchmarkOpeningPruningDepth4-4                1	4016534087 ns/op
    BenchmarkMidGamePruningDepth4-4                1	3511034720 ns/op
    PASS
    ok      github.com/polypmer/ghess	24.323s

After new check method

    BenchmarkSearchValid-4                     10000        290186 ns/op
    BenchmarkSearchValidSlow-4                 10000        305094 ns/op
    BenchmarkMidGamePruningDepth2-4               50      21060368 ns/op
    BenchmarkOpeningPruningDepth2-4              100      16300874 ns/op
    BenchmarkOpeningPruningDepth3-4               10     156043108 ns/op
    BenchmarkMidGamePruningDepth3-4               20     146742071 ns/op
    BenchmarkMidGameTwoPruningDepth3-4            20      60358058 ns/op
    BenchmarkOpeningPruningDepth4-4                1	2920717586 ns/op
    BenchmarkMidGamePruningDepth4-4                1	2233602005 ns/op
    PASS
    ok      github.com/polypmer/ghess	19.966s

After Clean Up Tests:

    BenchmarkMidGamePruningDepth2-4              100      25309438 ns/op
    BenchmarkOpeningPruningDepth2-4              100      17696422 ns/op
    BenchmarkOpeningPruningDepth3-4               10     183471772 ns/op
    BenchmarkMidGamePruningDepth3-4               10     147649582 ns/op
    BenchmarkMidGamePruningDepth3v2-4              2     623692681 ns/op
    BenchmarkOpeningPruningDepth4-4                1	2224223811 ns/op
    BenchmarkMidGamePruningDepth4-4                1	1748469646 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	5833172254 ns/op
    BenchmarkOpeningPruningDepth5-4                1	25067368045 ns/op
    BenchmarkMidGamePruningDepth5-4                1	12097247498 ns/op
    PASS
    ok      github.com/polypmer/ghess	56.751s


After Checking for Checkmate in Move Function:

    BenchmarkMidGamePruningDepth2-4              100      22388385 ns/op
    BenchmarkOpeningPruningDepth2-4              100      16043151 ns/op
    BenchmarkOpeningPruningDepth3-4                5     243147499 ns/op
    BenchmarkMidGamePruningDepth3-4                3     480370920 ns/op
    BenchmarkMidGamePruningDepth3v2-4              2     940100309 ns/op
    BenchmarkOpeningPruningDepth4-4                1	2995183075 ns/op
    BenchmarkMidGamePruningDepth4-4                1	2833983087 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	11650210041 ns/op
    BenchmarkOpeningPruningDepth5-4                1	45387620618 ns/op
    BenchmarkMidGamePruningDepth5-4                1	41308541742 ns/op
    PASS
    ok      github.com/polypmer/ghess	116.066s


With new checkCheck method:

    BenchmarkMidGamePruningDepth2-4              100      22448028 ns/op
    BenchmarkOpeningPruningDepth2-4              100      15045846 ns/op
    BenchmarkOpeningPruningDepth3-4                5     249524429 ns/op
    BenchmarkMidGamePruningDepth3-4                3     508006507 ns/op
    BenchmarkMidGamePruningDepth3v2-4              2     919852959 ns/op
    BenchmarkOpeningPruningDepth4-4                1	2828037344 ns/op
    BenchmarkMidGamePruningDepth4-4                1	2467769132 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	10528120071 ns/op
    BenchmarkOpeningPruningDepth5-4                1	44537287486 ns/op
    BenchmarkMidGamePruningDepth5-4                1	40094846903 ns/op
    PASS
    ok      github.com/polypmer/ghess	112.273s


Removing bytes to Upper in Favor of unicode ToLower (woah big gain)

    BenchmarkMidGamePruningDepth2-4              100      16925310 ns/op
    BenchmarkOpeningPruningDepth2-4              100      12206812 ns/op
    BenchmarkOpeningPruningDepth3-4               10     188508894 ns/op
    BenchmarkMidGamePruningDepth3-4                3     362414598 ns/op
    BenchmarkMidGamePruningDepth3v2-4              2     719300465 ns/op
    BenchmarkOpeningPruningDepth4-4                1	2205359134 ns/op
    BenchmarkMidGamePruningDepth4-4                1	2047525391 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	8160272334 ns/op
    BenchmarkOpeningPruningDepth5-4                1	35592436104 ns/op
    BenchmarkMidGamePruningDepth5-4                1	36422738754 ns/op
    PASS
    ok      github.com/polypmer/ghess	93.830s

Solving Mate in Three Puzzles

    BenchmarkMidGamePruningDepth2-4              100      17760837 ns/op
    BenchmarkOpeningPruningDepth2-4              100      11871520 ns/op
    BenchmarkOpeningPruningDepth3-4               10     189677602 ns/op
    BenchmarkMidGamePruningDepth3-4                3     423109193 ns/op
    BenchmarkMidGamePruningDepth3v2-4              2     763345004 ns/op
    BenchmarkOpeningPruningDepth4-4                1	2255497458 ns/op
    BenchmarkMidGamePruningDepth4-4                1	2156571491 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	8961648726 ns/op
    BenchmarkOpeningPruningDepth5-4                1	37199162918 ns/op
    BenchmarkMidGamePruningDepth5-4                1	39992108638 ns/op
    PASS
    ok      github.com/polypmer/ghess	100.381s


Fix isUpper method:

    BenchmarkMidGamePruningDepth2-4              100      18516595 ns/op
    BenchmarkOpeningPruningDepth2-4              100      12221270 ns/op
    BenchmarkOpeningPruningDepth3-4               10     192746436 ns/op
    BenchmarkMidGamePruningDepth3-4                3     396922504 ns/op
    BenchmarkMidGamePruningDepth3v2-4              2     914956021 ns/op
    BenchmarkOpeningPruningDepth4-4                1	2334247240 ns/op
    BenchmarkMidGamePruningDepth4-4                1	2202312773 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	8984800795 ns/op
    BenchmarkOpeningPruningDepth5-4                1	38791150950 ns/op
    BenchmarkMidGamePruningDepth5-4                1	38760424831 ns/op
    PASS
    ok      github.com/polypmer/ghess	101.368s

Redux SearchValid

    BenchmarkMidGamePruningDepth2-4              300       7301357 ns/op
    BenchmarkOpeningPruningDepth2-4              500       4243673 ns/op
    BenchmarkOpeningPruningDepth3-4               30      43071126 ns/op
    BenchmarkMidGamePruningDepth3-4               30      54478924 ns/op
    BenchmarkMidGamePruningDepth3v2-4             10     197960945 ns/op
    BenchmarkOpeningPruningDepth4-4                2     755648819 ns/op
    BenchmarkMidGamePruningDepth4-4                2     665499051 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	2091657012 ns/op
    BenchmarkOpeningPruningDepth5-4                1	8319886414 ns/op
    BenchmarkMidGamePruningDepth5-4                1	5527034454 ns/op
    PASS
    ok      github.com/polypmer/ghess	30.643s

Adding a Sixth Depth


    BenchmarkMidGamePruningDepth2-4              200       6498673 ns/op
    BenchmarkOpeningPruningDepth2-4              300       4322851 ns/op
    BenchmarkOpeningPruningDepth3-4               30      46424936 ns/op
    BenchmarkMidGamePruningDepth3-4               30      49564864 ns/op
    BenchmarkMidGamePruningDepth3v2-4              5     217738653 ns/op
    BenchmarkOpeningPruningDepth4-4                2     733054919 ns/op
    BenchmarkMidGamePruningDepth4-4                2     703439361 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	2104316956 ns/op
    BenchmarkOpeningPruningDepth5-4                1	8582348489 ns/op
    BenchmarkMidGamePruningDepth5-4                1	5890916097 ns/op
    BenchmarkMidGamePruningDepth5v2-4              1	75357024372 ns/op
    BenchmarkOpeningPruningDepth6-4                1	147948218393 ns/op
    BenchmarkMidGamePruningDepth6-4                1	176360843441 ns/op
    PASS
    ok      github.com/polypmer/ghess	428.519s

Slightly modify evaluation

    BenchmarkMidGamePruningDepth2-4              200       5758300 ns/op
    BenchmarkOpeningPruningDepth2-4              500       3567167 ns/op
    BenchmarkOpeningPruningDepth3-4               30      41444101 ns/op
    BenchmarkMidGamePruningDepth3-4               30      48787742 ns/op
    BenchmarkMidGamePruningDepth3v2-4             10     188809489 ns/op
    BenchmarkOpeningPruningDepth4-4                2     572857926 ns/op
    BenchmarkMidGamePruningDepth4-4                2     580117586 ns/op
    BenchmarkMidGamePruningDepth4v2-4              1	1721015023 ns/op
    BenchmarkOpeningPruningDepth5-4                1	7145903930 ns/op
    BenchmarkMidGamePruningDepth5-4                1	4277581657 ns/op
    BenchmarkMidGamePruningDepth5v2-4              1	66624534625 ns/op
    PASS
    ok      github.com/polypmer/ghess	93.219s

