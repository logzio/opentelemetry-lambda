# Ensure gem libs are on the load path in case environment hooks are ignored
require 'rubygems'

begin
  gem_root = ENV['GEM_PATH'] || ENV['GEM_HOME']
  if gem_root && Dir.exist?(gem_root)
    # Ensure RubyGems knows about our layer paths
    begin
      default_paths = Gem.default_path
      Gem.use_paths(ENV['GEM_HOME'] || gem_root, [gem_root, *default_paths].uniq)
      Gem.refresh
    rescue StandardError
      # ignore, we still amend $LOAD_PATH manually below
    end
    Dir.glob(File.join(gem_root, 'gems', '*', 'lib')).each do |dir|
      $LOAD_PATH.unshift(dir) unless $LOAD_PATH.include?(dir)
    end
  end
rescue StandardError
  # no-op: fall through to requires; errors will surface if libs are missing
end

require 'opentelemetry/sdk'
require 'opentelemetry/exporter/otlp'
require 'opentelemetry/instrumentation/all'

# We need to load the function code's dependencies, and _before_ any dependencies might
# be initialized outside of the function handler, bootstrap instrumentation.
def preload_function_dependencies
  default_task_location = '/var/task'

  handler_file = ENV.values_at('ORIG_HANDLER', '_HANDLER').compact.first&.split('.')&.first

  unless handler_file && File.exist?("#{default_task_location}/#{handler_file}.rb")
    OpenTelemetry.logger.warn { 'Could not find the original handler file to preload libraries.' }
    return nil
  end

  # Read as UTF-8 and scrub invalid bytes to avoid US-ASCII encoding errors
  source = File.read("#{default_task_location}/#{handler_file}.rb", mode: 'rb').force_encoding('UTF-8')
  source = source.sub(/^\uFEFF/, '') # strip UTF-8 BOM if present
  source = source.scrub
  libraries = source
                  .scan(/^\s*require\s+['"]([^'"]+)['"]/)
                  .flatten

  libraries.each do |lib|
    require lib
  rescue StandardError => e
    OpenTelemetry.logger.warn { "Could not load library #{lib}: #{e.message}" }
  end
  handler_file
end

handler_file = preload_function_dependencies

OpenTelemetry.logger.info { "Libraries in #{handler_file} have been preloaded." } if handler_file

OpenTelemetry::SDK.configure do |c|
  c.use_all()
end

def otel_wrapper(event:, context:)
  otel_wrapper = OpenTelemetry::Instrumentation::AwsLambda::Handler.new
  otel_wrapper.call_wrapped(event: event, context: context)
end
