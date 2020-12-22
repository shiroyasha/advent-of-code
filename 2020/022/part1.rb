# rubocop:disable all

require 'set'

class Player
  attr_reader :name
  attr_reader :deck

  def initialize(name, deck)
    @name = name
    @deck = deck
  end

  def draw
    @deck.shift
  end

  def add(card)
    @deck.push(card)
  end

  def has_cards?
    @deck.size > 0
  end

  def score
    @deck.map.with_index { |c, i| c * (@deck.size - i) }.inject(&:+)
  end

  def fingerprint
    @deck.join(",")
  end

  def to_s
    "#{@name}: #{@deck.inspect}"
  end
end

class Game2
  def initialize(player1, player2, id = 1)
    @id = id
    @player1 = player1
    @player2 = player2

    @visited1 = Set.new
    @visited2 = Set.new
    @round = 0
  end

  def run
    loop do
      @round += 1

      # puts ""
      # puts "--- Round #{@round} (Game #{@id}) ---"

      # puts @player1
      # puts @player2

      return @player2 unless @player1.has_cards?
      return @player1 unless @player2.has_cards?
      return @player1 if already_visited_state?

      record_state

      play_round
    end
  end

  def play_round
    c1 = @player1.draw
    c2 = @player2.draw

    # puts "#{@player1.name} draws #{c1}"
    # puts "#{@player2.name} draws #{c2}"

    winner = recursive_combat?(c1, c2) ? recursive_combat(c1, c2) : normal_combat(c1, c2)

    puts "#{winner.name} wins round #{@round} of game #{@id}"

    if winner.name == @player1.name
      @player1.add(c1)
      @player1.add(c2)
    else
      @player2.add(c2)
      @player2.add(c1)
    end
  end

  def recursive_combat(c1, c2)
    # puts "Recursive combat"

    p1 = Player.new(@player1.name, @player1.deck.clone[0...c1])
    p2 = Player.new(@player2.name, @player2.deck.clone[0...c2])

    Game2.new(p1, p2, @id + 1).run
  end

  def normal_combat(c1, c2)
    # puts "Normal combat"

    c1 > c2 ?  @player1 : @player2
  end

  def recursive_combat?(c1, c2)
    @player1.deck.size >= c1 && @player2.deck.size >= c2
  end

  def record_state
    @visited1.add(@player1.fingerprint)
    @visited2.add(@player2.fingerprint)
  end

  def already_visited_state?
    @visited1.member?(@player1.fingerprint) || @visited2.member?(@player2.fingerprint)
  end
end

module Day22
  def self.parse(path)
    input = File.read(path)

    p1, p2 = input.split("\n\n")

    d1 = p1.split("\n")[1..-1].map(&:to_i)
    d2 = p2.split("\n")[1..-1].map(&:to_i)

    [Player.new("Player 1", d1), Player.new("Player 2", d2)]
  end

  def self.play_one_round_simple(player1, player2)
    c1 = player1.draw
    c2 = player2.draw

    if c1 > c2
      player1.add(c1)
      player1.add(c2)
    else
      player2.add(c2)
      player2.add(c1)
    end
  end

  def self.play_game_simple(player1, player2)
    while player1.has_cards? && player2.has_cards?
      play_one_round_simple(player1, player2)
    end

    winner(player1, player2)
  end

  def self.calculate_score(player)
    player.score
  end

  def self.winner(player1, player2)
    if player1.has_cards?
      player1
    else
      player2
    end
  end

  def self.solve1(path)
    player1, player2 = parse(path)

    w = play_game_simple(player1, player2)

    calculate_score(w)
  end

  def self.solve2(path)
    player1, player2 = parse(path)

    w = Game2.new(player1, player2).run

    calculate_score(w)
  end
end

if __FILE__ == $0
  puts Day22.solve1(ARGV[0])
  puts Day22.solve2(ARGV[0])
end
