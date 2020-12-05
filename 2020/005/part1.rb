# rubocop:disable all

def binary_lookup(input, range)
  input.split("").each do |c|
    # puts "#{c} #{range.inspect}"

    half  = (range[1] - range[0]) / 2

    if c == "F" || c  == "L"
      range[1] = range[1] - half
    end

    if c == "B" || c == "R"
      range[0] = range[0] + half
    end

    # puts "=> #{range}"
  end

  range[0]
end

def process(input)
  row = binary_lookup(input[0..7], [0, 128])
  column = binary_lookup(input[7..-1], [0, 8])

  {row: row, column: column, id: row * 8 + column}
end

# p process("FBFBBFFRLR")
# p process("BFFFBBFRRR")
# p process("FFFBBBFRRR")
# p process("BBFFBBFRLL")

boarding_passes = File.read("input1.txt").split("\n").map { |i| process(i)[:id] }

(0...127).each do |row|
  missing = []

  (0...8).each do |column|
    id = row*8 + column

    if boarding_passes.include?(id)
      print "."
    else
      missing << [row, column, id]
      print "#"
    end
  end

  if missing.size == 1
    puts "  Here!!! #{missing[0]}"
  else
    puts
  end
end
