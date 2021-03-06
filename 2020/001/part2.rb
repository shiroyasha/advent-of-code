input = File.read("input.txt")

numbers = input.split("\n").map { |l| l.to_i }

def find(numbers)
  numbers.each do |n1|
    numbers.each do |n2|
      numbers.each do |n3|
        if n1 + n2 + n3 == 2020
          return [n1, n2, n3]
        end
      end
    end
  end
end

n1, n2, n3 = find(numbers)

puts(n1 * n2 * n3)
