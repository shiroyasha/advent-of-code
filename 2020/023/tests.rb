# rubocop:disable all

require_relative "./part2.rb"
require 'minitest/autorun'

class Day23Test < Minitest::Test

  # def test_one_round
  #   input = "389125467"

  #   assert_equal Day23.play(input, 0), "328915467"
  # end

  def test_10_rounds
    input = "389125467"

    res = Day23.solve1(input, 10)

    assert_equal res, 18
  end

  def test_100_rounds
    input = "186524973"

    res = Day23.solve1(input)

    assert_equal res, 20
  end

end
