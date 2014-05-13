require File.expand_path("../lib/weigh/version", __FILE__)

Gem::Specification.new do |gem|

  gem.name    = 'weigh'
  gem.version = Weigh::VERSION
  gem.date    = Date.today.to_s

  gem.summary     = "Weigh directory contents"
  gem.description = "Useful utility to weigh the directory contents and sort the output by size in a neatly formatted list."

  gem.author   = 'Zach Leslie'
  gem.email    = 'xaque208@gmail.com'
  gem.homepage = 'https://github.com/xaque208/weigh'

  # ensure the gem is built out of versioned files
   gem.files = Dir['Rakefile', '{bin,lib}/**/*', 'README*', 'LICENSE*'] & %x(git ls-files -z).split("\0")

   gem.executables << 'weigh'

end

