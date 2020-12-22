# rubocop:disable all

require_relative "./part1.rb"
require 'minitest/autorun'

class Day22Test < Minitest::Test

  def test_it_parse_the_input
    player1, player2 = Day22.parse("input0.txt")

    assert_equal player1.deck, [9, 2, 6, 3, 1]
    assert_equal player2.deck, [5, 8, 4, 7, 10]
  end

  def test_player_can_draw_a_card
    player = Player.new("P1", [1, 2, 3])

    assert_equal player.draw, 1
    assert_equal player.draw, 2
    assert_equal player.draw, 3
    assert_nil player.draw
  end

  def test_playing_one_round
    p1 = Player.new("P1", [1])
    p2 = Player.new("P2", [2])

    Day22.play_one_round_simple(p1, p2)

    assert_equal p1.deck, []
    assert_equal p2.deck, [2, 1]

    p1 = Player.new("P1", [9, 1])
    p2 = Player.new("P2", [2, 1])

    Day22.play_one_round_simple(p1, p2)

    assert_equal p1.deck, [1, 9, 2]
    assert_equal p2.deck, [1]
  end

  def test_playing_game
    player1, player2 = Day22.parse("input0.txt")

    Day22.play_game_simple(player1, player2)

    assert_equal player1.deck, []
    assert_equal player2.deck, [3, 2, 10, 6, 8, 5, 9, 4, 7, 1]
  end

  def test_calculating_the_score
    player1, player2 = Day22.parse("input0.txt")

    Day22.play_game_simple(player1, player2)

    assert_equal Day22.calculate_score(player2), 306
  end

  def test_infinite_loop_blocking
    player1 = Player.new("P1", [43, 19])
    player2 = Player.new("P2", [2, 29, 14])

    w = Game2.new(player1, player2).run

    assert_equal w.name, "P1"
  end

  def test_recursive_combat?
    player1 = Player.new("P1", [1, 2, 3])
    player2 = Player.new("P1", [1, 2, 3, 4])

    game = Game2.new(player1, player2)

    assert_equal game.recursive_combat?(4, 4), false
    assert_equal game.recursive_combat?(2, 2), true
    assert_equal game.recursive_combat?(6, 2), false
  end

  def test_game_normal_combat
    player1 = Player.new("P1", [1, 2, 3])
    player2 = Player.new("P2", [1, 2, 3, 4])

    game = Game2.new(player1, player2)

    assert_equal game.normal_combat(1, 2).name, "P2"
    assert_equal game.normal_combat(9, 2).name, "P1"
  end

  def test_playing_game_with_recurse
    player1, player2 = Day22.parse("input0.txt")

    winner = Game2.new(player1, player2).run

    assert_equal winner.name, "Player 2"
    assert_equal winner.score, 291
  end
end
