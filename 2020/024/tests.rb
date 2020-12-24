# rubocop:disable all

require_relative "./part1.rb"
require 'minitest/autorun'

class Day24Test < Minitest::Test

  def test_parse_line
    assert_equal Day24.parse_line("nwwswee"), ["nw", "w", "sw", "e", "e"]
  end

  def test_adjecent
    assert_equal Day24.adjecent_pos(Vec.new(0,0)), [
      Vec.new(1, 0),
      Vec.new(-1, 0),

      Vec.new(0, 1),
      Vec.new(1, 1),

      Vec.new(-1, -1),
      Vec.new(0, -1),
    ]
  end

  def test_find_coordinate
    assert_equal Day24.find_coordinate(["nw", "w", "sw", "e", "e"]), Vec.new(0, 0)
    assert_equal Day24.find_coordinate(["e", "se", "w"]), Day24.find_coordinate(["se"])
  end

  def test_calc_counts
    map = Set.new()
    map.add(Vec.new(0, 0))

    assert_equal Day24.calc_counts(map), {
      Vec.new(0, 0) => 0,
      Vec.new(1, 0) => 1,
      Vec.new(-1, 0) => 1,
      Vec.new(0, 1) => 1,
      Vec.new(1, 1) => 1,
      Vec.new(-1, -1) => 1,
      Vec.new(0, -1) => 1
    }
  end

  def test_it_can_solve_example_input
    assert_equal Day24.solve1("input0.txt"), 10
  end

  def test_it_can_solve2_on_example_input
    assert_equal Day24.solve2("input0.txt"), 2208
  end
end
