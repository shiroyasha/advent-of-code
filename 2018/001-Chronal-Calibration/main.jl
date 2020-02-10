function find(path::String)
  file  = open(path)
  lines = readlines(file)

  current_sum = 0
  visited = Set()

  while true
    for l in lines
      current_sum = current_sum + parse(Int, l)

      if in(current_sum, visited)
        return current_sum
      else
        push!(visited, current_sum)
      end
    end
  end
end

find("input1.txt") |> println
