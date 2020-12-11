# rubocop:disable all

input = File.read("input1.txt")

def display(matrix)
  matrix.each { |l| puts l.join("") }
end

def on_the_grid(matrix, i, j)
  return false if i < 0 || j < 0
  return false if i >= matrix.length
  return false if j >= matrix[0].length

  return true
end

def occupied?(matrix, i, j, incI, incJ)
  loop do
    i += incI
    j += incJ

    return false unless on_the_grid(matrix, i, j)

    return false if matrix[i][j] == "L"
    return true if matrix[i][j] == "#"
  end

  false
end

def transform(matrix)
  changes = 0

  new_matrix = matrix.map.with_index do |line, i|
    line.map.with_index do |c, j|
      next matrix[i][j] if matrix[i][j] == "."

      occupied = 0

      occupied = occupied + 1 if occupied?(matrix, i, j, -1, -1)
      occupied = occupied + 1 if occupied?(matrix, i, j, -1, 0)
      occupied = occupied + 1 if occupied?(matrix, i, j, -1, +1)

      occupied = occupied + 1 if occupied?(matrix, i, j, 0, -1)
      occupied = occupied + 1 if occupied?(matrix, i, j, 0, +1)

      occupied = occupied + 1 if occupied?(matrix, i, j, 1, -1)
      occupied = occupied + 1 if occupied?(matrix, i, j, 1, 0)
      occupied = occupied + 1 if occupied?(matrix, i, j, 1, +1)

      if matrix[i][j] == "L" && occupied == 0
        changes += 1
        '#'
      elsif matrix[i][j] == "#" && occupied >= 5
        changes += 1
        "L"
      else
        matrix[i][j]
      end
    end
  end

  [new_matrix, changes]
end

def occupied_total(matrix)
  matrix.flatten.count { |c| c == "#" }
end

def transform_until_stale(matrix)
  new_matrix, changes = transform(matrix)

  return new_matrix if changes == 0

  transform_until_stale(new_matrix)
end

matrix = input.split("\n").map { |l| l.split('') }

matrix = transform_until_stale(matrix)

puts "Part 2: #{occupied_total(matrix)}"
