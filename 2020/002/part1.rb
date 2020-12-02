Password = Struct.new(:from, :to, :character, :actual) do
  def valid?
    count = actual.split('').select { |c| c == character }.count

    count >= from and count <= to
  end
end

passwords = File.read("input.txt").split("\n").map do |line|
  limit, actual = line.split(": ")
  range, character = limit.split(" ")
  from, to = range.split("-")

  Password.new(from.to_i, to.to_i, character, actual)
end

puts passwords.select(&:valid?).count
