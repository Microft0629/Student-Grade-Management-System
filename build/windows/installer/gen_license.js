const fs = require('fs');

function esc(text) {
  let out = '';
  for (const ch of text) {
    const c = ch.codePointAt(0);
    if (c > 127) out += '\\u' + c + '?';
    else if (ch === '\n') out += '\\line\n';
    else if (ch === '\\' || ch === '{' || ch === '}') out += '\\' + ch;
    else out += ch;
  }
  return out;
}

const title   = esc('免责声明');
const body    = esc(`本软件（学生成绩管理系统）按"现状"提供，不提供任何形式的明示或暗示担保。\n\n使用本软件即表示您同意：\n\n1. 本软件仅供教育管理用途，不得用于任何非法目的。\n2. 开发者不对因使用本软件而导致的任何数据丢失或损坏承担责任。\n3. 用户应自行负责数据的备份与安全保管。\n4. 本软件所处理的学生成绩数据应遵守相关隐私法规。\n5. 开发者保留随时修改或终止软件服务的权利。`);
const footer  = esc('如不同意以上条款，请退出安装程序。');

const rtf = `{\\rtf1\\ansi\\deff0
{\\fonttbl{\\f0\\fnil\\fcharset134 SimSun;}{\\f1\\fnil\\fcharset134 Microsoft YaHei;}}
\\f0\\fs22

\\pard \\qc \\sa180 \\b \\fs26 ${title}\\b0 \\fs22 \\par

\\pard \\qj \\sa120 \\sl276 \\slmult1
${body}\\par

\\pard \\qc \\sa120
${footer}\\par
}`;

fs.writeFileSync('license.rtf', rtf);
console.log('Done');
