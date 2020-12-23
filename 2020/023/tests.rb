# rubocop:disable all

require_relative "./part1.rb"
require 'minitest/autorun'

class Day23Test < Minitest::Test

  # def test_one_round
  #   input = "389125467"

  #   assert_equal Day23.play(input, 0), "328915467"
  # end

  # def test_two_rounds
  #   input = "389125467"

  #   2.times.each do |i|
  #     input = Day23.play(input, i)
  #   end

  #   assert_equal input, "325467891"
  # end

  def test_10_rounds
    input = "389125467"

    res = Day23.play_all(input, 10)

    assert_equal res, "92658374"
  end

end
