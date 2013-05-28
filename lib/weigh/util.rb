module Weigh
  module Util

    def self.neat_size(bytes)
      # return a human readble size from the bytes supplied
      bytes = bytes.to_f
      if bytes > 2 ** 40       # TiB: 1024 GiB
        neat = sprintf("%.2f TiB", bytes / 2**40)
      elsif bytes > 2 ** 30    # GiB: 1024 MiB
        neat = sprintf("%.2f GiB", bytes / 2**30)
      elsif  bytes > 2 ** 20   # MiB: 1024 KiB
        neat = sprintf("%.2f MiB", bytes / 2**20)
      elsif bytes > 2 ** 10    # KiB: 1024 B
        neat = sprintf("%.2f KiB", bytes / 2**10)
      else                    # bytes
        neat = sprintf("%.0f bytes", bytes)
      end
      neat
    end

    def self.sum_dir(dir,verbose=false)
      # return the size of a given directory
      #"Entering: #{dir}"
      count=0
      dir_size=0
      data={}
      Find.find(dir) do |path|
        count += 1
        next if FileTest.symlink?(path)
        next if dir == path
        if FileTest.directory?(path)
          ret = sum_dir(path,verbose)
          size = ret[:dir_size]
          count += ret[:count]
          dir_size += size
          Find.prune
        else
          size = FileTest.size(path)
          #puts "File: #{path} is #{size}"
          puts "Found zero size file: #{path}" if verbose
          dir_size += size
        end
      end
      #puts "Exiting: #{dir} with #{dir_size}"
      data[:dir_size] = dir_size
      data[:count] = count
      data
    end

    def self.report(summary,total_size)
      summary.sort{|a,b| a[1]<=>b[1]}.each { |elem|
        size     = elem[1]
        filename = elem[0]
        puts sprintf("%15s   %s\n", neat_size(size), filename)
      }

      puts sprintf("%16s %s\n", "---", "---")
      puts sprintf("%15s   %s\n", self.neat_size(total_size), ":total size")
      puts sprintf("%16s %s\n", "---", "---")
    end

    def sum_dir(dir,verbose=false)
      # return the size of a given directory
      #"Entering: #{dir}"
      count=0
      dir_size=0
      data={}
      Find.find(dir) do |path|
        count += 1
        next if FileTest.symlink?(path)
        next if dir == path
        if FileTest.directory?(path)
          ret = sum_dir(path,verbose)
          size = ret[:dir_size]
          count += ret[:count]
          dir_size += size
          Find.prune
        else
          size = FileTest.size(path)
          #puts "File: #{path} is #{size}"
          puts "Found zero size file: #{path}" if verbose
          dir_size += size
        end
      end
      #puts "Exiting: #{dir} with #{dir_size}"
      data[:dir_size] = dir_size
      data[:count] = count
      data
    end
  end
end
