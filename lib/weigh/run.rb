require 'weigh/util'
require 'find'

module Weigh
  class Run

    attr_accessor :verbose, :pathlist

    def initialize(pathlist=['.'])
      @pathlist = pathlist

      @data              = {
        :total_size => 0,
        :summary    => {},
        :count      => 0
      }
    end

    def verbose=(value)
      @verbose = value
    end

    def verbose
      @verbose
    end

    def pathlist=(value)
      @pathlist = value
    end

    def pathlist
      @pathlist
    end

    def run
      # Dig into each path given
      @pathlist.each do |p|
        Find.find(p) do |path|
          @data[:count] += 1

          # Skip symlinks
          next if FileTest.symlink?(path)

          if FileTest.directory?(path)

            # Skip the path that we are already on
            next if p == path

            # Summarize the directory data
            ret = Weigh::Util.sum_dir(path,@verbose)
            dir_size = ret[:dir_size]

            # Skip empty directories
            next if dir_size == 0

            # Record the data from the directory
            @data[:count]      += ret[:count]
            @data[:total_size] += dir_size

            # Add a trailing slash to the key for directories
            pathname = path + "/"

            @data[:summary][pathname] = dir_size

            Find.prune
          else
            # Store the size of the current file
            size = FileTest.size(path)

            # Don't count zero sized files
            next if size == 0

            @data[:total_size] += size
            @data[:summary][path] = size
          end
        end
      end

      @data
    end

    def report
      Weigh::Util.report @data
    end
  end
end
