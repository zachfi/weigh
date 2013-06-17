require 'weigh/util'
require 'pp'
require 'find'

module Weigh
  class Runner

    attr_reader :flags

    def initialize(flags)
      @flags = flags
    end

    def run
      data = {}
      total_size = 0
      summary = {}
      count = 0
      sumdep = @flags.depth

      @flags.pathlist.each do |p|
        curdep = 0
        Find.find(p) do |path|
          count += 1
          next if FileTest.symlink?(path)
          if FileTest.directory?(path)
            next if p == path
            ret = Weigh::Util.sum_dir(path,@flags.verbose)
            dir_size = ret[:dir_size]
            next if dir_size == 0
            count += ret[:count]
            total_size += dir_size
            summary["#{path}/"] = dir_size
            Find.prune
          else
            size = FileTest.size(path)
            next if size == 0
            total_size += size
            summary["#{path}"] = size
          end
        end
      end
      data[:summary]    = summary
      data[:total_size] = total_size

      data
    end
  end
end
