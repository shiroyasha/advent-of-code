# rubocop:disable all

require_relative "./part1.rb"
require 'minitest/autorun'

class Day21Test < Minitest::Test

  def test_it_can_solve_example_input
    assert_equal Day21.solve1("input0.txt"), 5
  end
end
