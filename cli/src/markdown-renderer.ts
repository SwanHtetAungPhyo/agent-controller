import chalk from 'chalk';

export function renderMarkdown(text: string): string {
  let rendered = text;

  // Headers
  rendered = rendered.replace(/^### (.*$)/gim, chalk.bold.cyan('### $1'));
  rendered = rendered.replace(/^## (.*$)/gim, chalk.bold.magenta('## $1'));
  rendered = rendered.replace(/^# (.*$)/gim, chalk.bold.yellow('# $1'));

  // Bold
  rendered = rendered.replace(/\*\*(.*?)\*\*/g, chalk.bold.white('$1'));
  
  // Italic
  rendered = rendered.replace(/\*(.*?)\*/g, chalk.italic('$1'));
  
  // Code blocks
  rendered = rendered.replace(/```([\s\S]*?)```/g, (match, code) => {
    return '\n' + chalk.bgBlack.yellow(code.trim()) + '\n';
  });
  
  // Inline code
  rendered = rendered.replace(/`(.*?)`/g, chalk.yellow('$1'));
  
  // Links
  rendered = rendered.replace(/\[(.*?)\]\((.*?)\)/g, chalk.blue.underline('$1') + chalk.gray(' ($2)'));
  
  // Bullet points
  rendered = rendered.replace(/^- (.*$)/gim, chalk.cyan('  • ') + '$1');
  rendered = rendered.replace(/^\* (.*$)/gim, chalk.cyan('  • ') + '$1');
  
  // Numbers
  rendered = rendered.replace(/\$([0-9,]+\.?[0-9]*)/g, chalk.green.bold('$$$1'));
  rendered = rendered.replace(/([0-9]+\.?[0-9]*)%/g, chalk.cyan.bold('$1%'));
  
  // Blockquotes
  rendered = rendered.replace(/^> (.*$)/gim, chalk.gray.italic('│ $1'));
  
  return rendered;
}

export async function animateText(text: string, delay: number = 3): Promise<void> {
  for (let i = 0; i < text.length; i++) {
    process.stdout.write(text[i]);
    if (delay > 0) {
      await new Promise(resolve => setTimeout(resolve, delay));
    }
  }
}

export async function animateMarkdown(text: string, lineDelay: number = 20): Promise<void> {
  const rendered = renderMarkdown(text);
  const lines = rendered.split('\n');
  
  for (const line of lines) {
    process.stdout.write(line + '\n');
    if (lineDelay > 0) {
      await new Promise(resolve => setTimeout(resolve, lineDelay));
    }
  }
}
