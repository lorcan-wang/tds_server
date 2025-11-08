module.exports = function (api) {
  api.cache(true);
  return {
    presets: ['babel-preset-expo'],
    plugins: [
      [
        'module-resolver',
        {
          extensions: ['.tsx', '.ts', '.js', '.json'],
          alias: {
            '@api': './src/api',
            '@components': './src/components',
            '@navigation': './src/navigation',
            '@screens': './src/screens',
            '@store': './src/store',
            '@types': './src/types',
            '@utils': './src/utils',
            '@providers': './src/providers'
          }
        }
      ],
      'react-native-reanimated/plugin'
    ]
  };
};
