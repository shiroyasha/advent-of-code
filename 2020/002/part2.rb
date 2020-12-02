Password = Struct.new(:pos1, :pos2, :character, :actual) do
  def valid?
    characters = actual.split('')

    (characters[pos1 - 1] == character) ^ (characters[pos2 - 1] == character)
  end
end

passwords = File.read('input.txt').split("\n").map do |line|
  limit, actual = line.split(': ')
  range, character = limit.split(' ')
  from, to = range.split('-')

  Password.new(from.to_i, to.to_i, character, actual)
end

puts passwords.select(&:valid?).count
