import fs from 'fs';
import { glob } from 'glob';
import path from 'path';
import { fileURLToPath } from 'url';

// ESM下__dirname替代方案
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Regex for extracting i18n keys - only matches $t("xxx") format
const I18N_PATTERNS = [
  /\$t\(\s*['"`]([^'"`]+)['"`]\s*\)/g, // $t('key') or $t("key")
];

// Extract i18n keys from file content
function extractKeysFromContent(content, filePath) {
  const keys = new Set();

  I18N_PATTERNS.forEach(pattern => {
    let match;
    pattern.lastIndex = 0;
    while ((match = pattern.exec(content)) !== null) {
      const key = match[1];
      // Filter out dynamic keys and strings that are obviously not translation keys
      if (
        key &&
        !key.includes('+') && // Filter dynamic concatenation
        !key.includes('${') && // Filter template strings
        !key.includes('/') && // Filter paths
        !key.includes('@') && // Filter Vue paths
        !key.includes('.vue') && // Filter filenames
        !key.includes('http') && // Filter URLs
        !key.includes('api/') && // Filter API paths
        !key.includes('#') && // Filter CSS selectors/anchors
        !key.includes('<') && // Filter HTML
        !key.includes('>') && // Filter HTML
        !key.includes('=') && // Filter assignment expressions
        !key.includes('&') && // Filter URL params
        key.length > 1 && // Filter single character
        key.length < 200 && // Filter overly long strings
        !key.match(/^[.]+$/) && // Filter only dots
        !key.match(/^[0-9]+$/) && // Filter pure numbers
        !key.match(/^[a-zA-Z]:[\\\/]/) && // Filter Windows paths
        !key.match(/^\w+:\/\//) // Filter protocol URLs
      ) {
        keys.add(key);
      }
    }
  });

  return Array.from(keys);
}

// Extract keys from file
function extractKeysFromFile(filePath) {
  try {
    const content = fs.readFileSync(filePath, 'utf-8');
    return extractKeysFromContent(content, filePath);
  } catch (e) {
    console.warn(`Failed to read file: ${filePath}`, e.message);
    return [];
  }
}

// Read language file - supports ES Module format
async function readLanguageFile(filePath) {
  try {
    if (!fs.existsSync(filePath)) {
      console.warn(`Language file does not exist: ${filePath}`);
      return {};
    }

    // Read file content and parse
    const content = fs.readFileSync(filePath, 'utf-8');

    // Handle ES Module format: export default { ... }
    const exportMatch = content.match(/export\s+default\s+(\{[\s\S]*\})/);
    if (exportMatch) {
      try {
        // Use eval to parse object (assume language file is safe)
        const objStr = exportMatch[1];
        const langData = eval(`(${objStr})`);
        return cleanupNestedStructure(langData);
      } catch (e) {
        console.warn(`Failed to parse language file: ${filePath}`, e.message);
        return {};
      }
    }

    // If not ES Module, try dynamic import
    let langData;
    try {
      const imported = await import(pathToFileUrl(filePath));
      langData = imported.default || imported;
    } catch (e) {
      console.warn(`Failed to import language file: ${filePath}`, e.message);
      return {};
    }
    return cleanupNestedStructure(langData);
  } catch (e) {
    console.warn(`Failed to read language file: ${filePath}`, e.message);
    return {};
  }
}

function pathToFileUrl(filePath) {
  let absolutePath = filePath;
  if (!path.isAbsolute(filePath)) {
    absolutePath = path.resolve(process.cwd(), filePath);
  }
  let url = 'file://' + absolutePath;
  // Windows fix
  if (process.platform === 'win32') {
    url = url.replace(/\\/g, '/');
  }
  return url;
}

// Get all keys from object (supports nested)
function getAllKeys(obj, prefix = '') {
  const keys = [];
  for (const key in obj) {
    const fullKey = prefix ? `${prefix}.${key}` : key;
    const value = obj[key];
    if (
      value &&
      typeof value === 'object' &&
      !Array.isArray(value) &&
      Object.keys(value).length > 0
    ) {
      const subKeys = Object.keys(value);
      const hasValidSubKeys = subKeys.some(subKey => subKey.trim() !== '');
      if (hasValidSubKeys) {
        keys.push(...getAllKeys(value, fullKey));
      } else {
        keys.push(fullKey);
      }
    } else {
      keys.push(fullKey);
    }
  }
  return keys;
}

// Set nested value in object
function setNestedValue(obj, key, value) {
  const isPotentialNestedKey =
    key.includes('.') &&
    !key.startsWith('.') &&
    !key.endsWith('.') &&
    !/[,，。？?!！；;：:""''""''()（）\[\]【】{}｛｝<>《》]/.test(key) &&
    key
      .split('.')
      .every(part => part.trim() !== '' && /^[a-zA-Z_$][a-zA-Z0-9_$]*$/.test(part.trim()));
  if (isPotentialNestedKey) {
    const keys = key.split('.');
    let current = obj;
    for (let i = 0; i < keys.length - 1; i++) {
      const currentKey = keys[i];
      if (
        !current[currentKey] ||
        typeof current[currentKey] !== 'object' ||
        Array.isArray(current[currentKey])
      ) {
        current[currentKey] = {};
      }
      current = current[currentKey];
    }
    current[keys[keys.length - 1]] = value;
  } else {
    obj[key] = value;
  }
}

// Clean up incorrect nested structure
function cleanupNestedStructure(obj) {
  const cleaned = {};
  function traverse(current, target, prefix = '') {
    for (const key in current) {
      const value = current[key];
      const fullKey = prefix ? `${prefix}.${key}` : key;
      if (value && typeof value === 'object' && !Array.isArray(value)) {
        const subKeys = Object.keys(value);
        if (subKeys.length === 1 && subKeys[0].trim() === '') {
          target[key] = value[''];
          console.log(`🔧 Fixed incorrect nesting: "${key}" -> "${value['']}" `);
        } else if (subKeys.length === 0) {
          target[key] = key;
        } else if (subKeys.every(subKey => subKey.trim() === '')) {
          target[key] = key;
          console.log(`🔧 Fixed empty key nesting: "${key}"`);
        } else {
          const hasValidSubKeys = subKeys.some(subKey => subKey.trim() !== '');
          if (hasValidSubKeys) {
            if (!target[key]) target[key] = {};
            traverse(value, target[key], fullKey);
          } else {
            target[key] = key;
          }
        }
      } else {
        target[key] = value;
      }
    }
  }
  traverse(obj, cleaned);
  return cleaned;
}

// 新增：严格按模板对象结构和顺序生成目标语言对象
function buildLangByTemplate(templateObj, langObj) {
  if (typeof templateObj !== 'object' || templateObj === null) return templateObj;
  const result = Array.isArray(templateObj) ? [] : {};
  for (const key of Object.keys(templateObj)) {
    if (typeof templateObj[key] === 'object' && templateObj[key] !== null) {
      result[key] = buildLangByTemplate(templateObj[key], langObj?.[key] || {});
    } else {
      result[key] = langObj && langObj[key] !== undefined ? langObj[key] : key;
    }
  }
  return result;
}

// Main function
async function main() {
  console.log('Start extracting i18n keys...');
  const vueGlobs = ['src/**/*.vue', 'src/**/*.js', 'src/**/*.ts'];
  let sourceFiles = [];
  vueGlobs.forEach(pattern => {
    const files = glob.sync(pattern, { cwd: process.cwd(), absolute: true });
    console.log(`Pattern: ${pattern}, matched: ${files.length}`);
    sourceFiles = sourceFiles.concat(files);
  });
  sourceFiles = sourceFiles.filter(f => {
    try {
      return fs.existsSync(f) && fs.statSync(f).isFile() && fs.statSync(f).size > 0;
    } catch (e) {
      return false;
    }
  });
  console.log(`Total files scanned: ${sourceFiles.length}`);
  const allKeys = new Set();
  sourceFiles.forEach(file => {
    const keys = extractKeysFromFile(file);
    keys.forEach(key => allKeys.add(key));
  });
  console.log(`Extracted keys: ${allKeys.size}`);
  const keyArray = Array.from(allKeys);
  console.log('Sample keys:');
  keyArray.slice(0, 20).forEach(key => console.log(`  - ${key}`));
  if (keyArray.length > 20) {
    console.log(`  ... and ${keyArray.length - 20} more keys`);
  }
  const keysWithDots = keyArray.filter(key => key.includes('.'));
  if (keysWithDots.length > 0) {
    console.log(`\nKeys with dot (${keysWithDots.length}):`);
    keysWithDots.slice(0, 10).forEach(key => console.log(`  - "${key}"`));
    if (keysWithDots.length > 10) {
      console.log(`  ... and ${keysWithDots.length - 10} more keys with dot`);
    }
  }
  const langFiles = [
    path.join(process.cwd(), 'src/lang/zh-CN.js'),
    path.join(process.cwd(), 'src/lang/en-US.js'),
    path.join(process.cwd(), 'src/lang/zh-TW.js'),
  ];
  const langData = {};
  for (const file of langFiles) {
    const lang = path.basename(file, '.js');
    langData[lang] = await readLanguageFile(file);
    console.log(`Read language file ${lang}: ${Object.keys(langData[lang]).length} keys`);
    const originalKeyCount = Object.keys(langData[lang]).length;
    langData[lang] = cleanupNestedStructure(langData[lang]);
    const cleanedKeyCount = Object.keys(langData[lang]).length;
    if (originalKeyCount !== cleanedKeyCount) {
      console.log(`🔧 ${lang} after cleanup: ${cleanedKeyCount} keys`);
    }
  }
  // 以 en-US.js 为模板，严格对齐
  const enTemplate = langData['en-US'];
  const zhCNAligned = buildLangByTemplate(enTemplate, langData['zh-CN']);
  const zhTWAligned = buildLangByTemplate(enTemplate, langData['zh-TW']);
  // 写回 zh-CN.js
  fs.writeFileSync(
    path.join(process.cwd(), 'src/lang/zh-CN.js'),
    `export default ${JSON.stringify(zhCNAligned, null, 2)};`,
    'utf-8'
  );
  // 写回 zh-TW.js
  fs.writeFileSync(
    path.join(process.cwd(), 'src/lang/zh-TW.js'),
    `export default ${JSON.stringify(zhTWAligned, null, 2)};`,
    'utf-8'
  );
  // en-US.js 只保留自身结构
  fs.writeFileSync(
    path.join(process.cwd(), 'src/lang/en-US.js'),
    `export default ${JSON.stringify(enTemplate, null, 2)};`,
    'utf-8'
  );
  console.log('\n所有语言文件已严格按 en-US.js 对齐，key 顺序、结构、行数完全一致。');
  console.log('\nExtraction complete!');
}

if (import.meta.url === process.argv[1] || import.meta.url === `file://${process.argv[1]}`) {
  main();
}

export {
  extractKeysFromContent,
  extractKeysFromFile,
  getAllKeys,
  readLanguageFile,
  setNestedValue,
};
