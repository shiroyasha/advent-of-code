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

    assert_equal tiles[0].edges.size, 4
    assert_equal tiles[0].edges[0], "..##.#..#."
    assert_equal tiles[0].edges[1], "...#.##..#"
    assert_equal tiles[0].edges[2], "..###..###"
    assert_equal tiles[0].edges[3], ".#####..#."
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

  def test_solve_simple
    tiles = Day20.load("input0.txt").select { |t| [1951, 2311].include?(t.id) }

    s = Solver.new(tiles)

    assert_equal s.solve, [
      [1951, 2311]
    ]
  end

end
