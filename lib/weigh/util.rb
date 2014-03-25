module Weigh
  module Util

    # Convert a byte count into something a bit more human friendly
    #
    # @param [Int] a byte count
    #
    # @return [String] human readable byte count 
    def self.neat_size(bytes)
      # return a human readble size from the bytes supplied
      bytes = bytes.to_f
      if bytes > 2 ** 50       # PiB: 1024 TiB
        neat = sprintf("%.2f PiB", bytes / 2**50)
      elsif bytes > 2 ** 40    # TiB: 1024 GiB
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

    # Sumarize the size of the given directory
    #
    # @param [String] full path to directory
    # @param [Boolean] increase verbosity
    #
    # @return Hash {:dir_size => Hash, :count => Int}
    def self.sum_dir(dir,verbose=false)
      # return the size of a given directory
      #"Entering: #{dir}"
      count    = 0
      dir_size = 0
      data     = {}

      Find.find(dir) do |path|
        begin
          puts path if verbose
          count += 1
          if FileTest.symlink?(path)
            puts "skipping symlink " + path if verbose
            next
          end
          if dir == path and File.directory?(path)
            puts "skipping current directory " + path if verbose
            next
          end
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
        rescue IOError
          puts "file vanished: " + path
        end
      end
      data[:dir_size] = dir_size
      data[:count] = count
      data
    end

    # Print a summary of the received data
    #
    # @param [Hash]
    #
    # @returns [Hash] the unmodified data
    def self.report(data)
      data[:summary].sort{|a,b| a[1]<=>b[1]}.each { |elem|
        size     = elem[1]
        filename = elem[0]
        puts sprintf("%15s   %s\n", neat_size(size), filename)
      }

      puts sprintf("%16s %s\n", "---", "---")
      puts sprintf("%15s   %s\n", self.neat_size(data[:total_size]), ":total size")
      puts sprintf("%16s %s\n", "---", "---")

      data
    end
  end
end
