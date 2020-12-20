# rubocop:disable all

require_relative "./part1.rb"

require 'minitest/autorun'

class Day20Test < Minitest::Test

  def test_parsing_of_a_tile
    tiles = Day20.load("input0.txt")

    assert_equal tiles.size, 9
    assert_equal tiles[0].id, 2311
    assert_equal tiles[0].map.size, 10
    assert_equal tiles[0].map[0].size, 10
  end

  def test_tiles_have_edges
    tiles = Day20.load("input0.txt")

    assert_equal tiles[0].up, "..##.#..#."
    assert_equal tiles[0].down, "..###..###"

    assert_equal tiles[0].left, ".#####..#."
    assert_equal tiles[0].right, "...#.##..#"
  end

  def test_tile_can_be_rotated
    map = [
      "..##.#..#.",
      "##..#.....",
      "#...##..#.",
      "####.#...#",
      "##.##.###.",
      "##...#.###",
      ".#.#.#..##",
      "..#....#..",
      "###...#.#.",
      "..###..###"
    ]

    t = Tile.new(1, map)

    assert_equal t.rotate.map, [
      ".#..#####.",
      ".#.####.#.",
      "###...#..#",
      "#..#.##..#",
      "#....#.##.",
      "...##.##.#",
      ".#...#....",
      "#.#.##....",
      "##.###.#.#",
      "#..##.#..."
    ]
  end

  def test_tile_can_be_flipped
    map = [
      "..##.#..#.",
      "##..#.....",
      "#...##..#.",
      "####.#...#",
      "##.##.###.",
      "##...#.###",
      ".#.#.#..##",
      "..#....#..",
      "###...#.#.",
      "..###..###"
    ]

    t = Tile.new(1, map)

    assert_equal t.flip.map, [
      ".#..#.##..",
      ".....#..##",
      ".#..##...#",
      "#...#.####",
      ".###.##.##",
      "###.#...##",
      "##..#.#.#.",
      "..#....#..",
      ".#.#...###",
      "###..###.."
    ]
  end

  def test_desk_edges
    tiles = Day20.load("input0.txt")

    desk = Desk.new(3)
    desk.set(0, 0, tiles[0])

    assert_equal desk.edges, Set.new([
      [-1, 0], [1, 0], [0, -1], [0, 1],
    ])

    desk.set(1, 0, tiles[1])

    assert_equal desk.edges, Set.new([
      [-1, 0], [0, -1], [0, 1], [2, 0], [1, -1], [1, 1]
    ])
  end

  def test_desk_can_fit?
    tiles = Day20.load("input0.txt").select { |t| [1951, 2311].include?(t.id) }

    desk = Desk.new(3)
    desk.set(0, 0, tiles[1])

    assert_equal desk.can_fit?(1, 0, tiles[0]), true
  end

  def test_solve_simple
    tiles = Day20.load("input0.txt").select { |t| [1951, 2311].include?(t.id) }

    s = Solver.new(tiles)

    assert_equal s.solve.as_map, [
      [1951, 2311]
    ]
  end

  def test_solve_three
    tiles = Day20.load("input0.txt").select { |t| [1951, 2311, 3079].include?(t.id) }

    s = Solver.new(tiles)

    assert_equal s.solve.as_map, [
      [1951, 2311, 3079]
    ]
  end

  def test_solve_full
    tiles = Day20.load("input0.txt")

    s = Solver.new(tiles)

    assert_equal s.solve.as_map, [
      [1951, 2311, 3079], [2729, 1427, 2473], [2971, 1489, 1171]
    ]
  end

  def test_desk_to_image
    tiles = Day20.load("input0.txt")

    s = Solver.new(tiles)

    assert_equal s.solve.image, [
      ".#.#..#.##...#.##..#####",
      "###....#.#....#..#......",
      "##.##.###.#.#..######...",
      "###.#####...#.#####.#..#",
      "##.#....#.##.####...#.##",
      "...########.#....#####.#",
      "....#..#...##..#.#.###..",
      ".####...#..#.....#......",
      "#..#.##..#..###.#.##....",
      "#.####..#.####.#.#.###..",
      "###.#.#...#.######.#..##",
      "#.####....##..########.#",
      "##..##.#...#...#.#.#.#..",
      "...#..#..#.#.##..###.###",
      ".#.#....#.##.#...###.##.",
      "###.#...#..#.##.######..",
      ".#.#.###.##.##.#..#.##..",
      ".####.###.#...###.#..#.#",
      "..#.#..#..#.#.#.####.###",
      "#..####...#.#.#.###.###.",
      "#####..#####...###....##",
      "#.##..#..#...#..####...#",
      ".#.###..##..##..####.##.",
      "...###...##...#...#..###",
    ]
  end

  def test_is_monster
    image = [
      "            .     # ",
      "#. . ##  . ##    ###",
      " #  #  #  #  #  #   "
    ]

    s = SearchForMonstrer.new(image)
    assert s.is_monster(0, 0)

    image = [
      "            .     # ",
      ".. . ##  . ##    ###",
      " #  #  #  #  #  #   "
    ]

    s = SearchForMonstrer.new(image)
    refute s.is_monster(0, 0)
  end

end
