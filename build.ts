// build.ts - Generado por la ia DeepSeek con un promp un poco especifico genero un maravilloso codigo. - las mejoras tamien las realizo con promps cortitos. - 
import { existsSync, mkdirSync, readdirSync, statSync, copyFileSync, rmSync, watch } from 'fs';
import { join, relative, dirname, extname, basename } from 'path';
import { spawn, execSync, ChildProcess } from 'child_process';

const config = {
    tsInput: './internal/pages/ts',
    jsOutput: './public/js',
    jsBuild: './build/public/js',
    tailwindInput: './internal/pages/css',
    tailwindOutput: './public/css',
    tailwindBuild: './build/public/css',
    goMain: './main.go',
    goInput: './internal',
    goOutput: './',
    goBuild: './build',
    publicInput: './public',
    publicIgnore: ['js', 'css'],
    publicBuild: './build/public',
};

const isWatchMode = process.argv.includes('--watch');
const isBuildMode = !isWatchMode;

async function buildTypeScript(src: string, dest: string) {
    if (!existsSync(src)) return;
    
    if (!existsSync(dest)) {
        mkdirSync(dest, { recursive: true });
    }

    const files = readdirSync(src);
    for (const file of files) {
        const srcPath = join(src, file);
        const destPath = join(dest, file.replace('.ts', '.js'));

        if (statSync(srcPath).isDirectory()) {
            await buildTypeScript(srcPath, join(dest, file));
        } else if (file.endsWith('.ts')) {
            const relativePath = relative(process.cwd(), srcPath);
            const result = Bun.build({
                entrypoints: [relativePath],
                outdir: dirname(destPath),
                target: 'browser',
                format: 'esm',
                sourcemap: isWatchMode ? 'inline' : 'none',
                minify: isBuildMode,
            });
            
            if (!(await result).success) {
                console.error('TypeScript build failed');
                process.exit(1);
            }
        }
    }
}

async function buildTailwind(src: string, dest: string) {
    if (!existsSync(src)) return;

    if (!existsSync(dest)) {
        mkdirSync(dest, { recursive: true });
    }

    const inputFile = join(src, 'main.css');
    const outputFile = join(dest, 'main.css');

    try {
        execSync(`bunx --bun tailwindcss -i ${inputFile} -o ${outputFile} ${
        isBuildMode ? '--minify' : ''
        }`, { stdio: 'inherit' });
    } catch (error) {
        console.error('Tailwind build failed:', error);
        process.exit(1);
    }
}

function copyPublic(src: string, dest: string, ignore: string[]) {
    if (!existsSync(src)) return;

    if (!existsSync(dest)) {
        mkdirSync(dest, { recursive: true });
    }

    const items = readdirSync(src);
    for (const item of items) {
        if (ignore.includes(item)) continue;

        const srcPath = join(src, item);
        const destPath = join(dest, item);

        if (statSync(srcPath).isDirectory()) {
            copyPublic(srcPath, destPath, ignore);
        } else {
            copyFileSync(srcPath, destPath);
        }
    }
}

function runGoDev() {
    const goProcess = spawn('go', ['run', config.goMain], {
        stdio: 'inherit',
        cwd: config.goOutput,
    });

    goProcess.on('error', (error) => {
        console.error('Go execution failed:', error);
        process.exit(1);
    });

    return goProcess;
}

async function buildGo() {
    try {
        execSync(`go build -o ${config.goBuild} ${config.goMain}`, {
        stdio: 'inherit',
        });
    } catch (error) {
        console.error('Go build failed:', error);
        process.exit(1);
    }
}

// Función para encontrar todos los archivos .go en un directorio
function findGoFiles(dir: string): string[] {
    if (!existsSync(dir)) return [];
    
    const files: string[] = [];
    const items = readdirSync(dir);
  
    for (const item of items) {
        const fullPath = join(dir, item);
        if (statSync(fullPath).isDirectory()) {
            files.push(...findGoFiles(fullPath));
        } else if (item.endsWith('.go')) {
            files.push(fullPath);
        }
    }
    
    return files;
}

// async function stopGoProcess(goProcess: ChildProcess): Promise<void> {
//     return new Promise((resolve) => {
//         if (!goProcess.pid) {
//             resolve();
//             return;
//         }

//         // Enviar señal SIGINT (equivalente a Ctrl+C)
//         process.kill(goProcess.pid, 'SIGINT');
//         // process.kill(goProcess.pid, 'SIGTERM');
        
//         // Esperar a que el proceso termine gracefulmente
//         const timeout = setTimeout(() => {
//             // Si no termina después de 3 segundos, forzar la terminación
//             goProcess.kill('SIGKILL');
//             resolve();
//         }, 3000);

//         goProcess.once('exit', () => {
//             clearTimeout(timeout);
//             resolve();
//         });
//     });
// }

async function stopGoProcess(goProcess: ChildProcess): Promise<void> {
    return new Promise((resolve) => {
        if (!goProcess.pid || goProcess.exitCode !== null) {
            resolve();
            return;
        }

        try {
            // Enviar señal SIGTERM
            process.kill(goProcess.pid, 'SIGTERM');
        } catch (err: any) {
            if (err.code === 'ESRCH') {
                // Proceso ya no existe, lo damos por terminado
                resolve();
                return;
            }
            throw err; // otro error inesperado
        }

        // Esperar a que el proceso termine graceful
        const timeout = setTimeout(() => {
            if (goProcess.exitCode === null) {
                goProcess.kill('SIGKILL');
            }
            resolve();
        }, 9000);

        goProcess.once('exit', () => {
            clearTimeout(timeout);
            resolve();
        });
    });
}

async function main() {
    // Clean build directory
    if (isBuildMode && existsSync(config.goBuild)) {
        rmSync(config.goBuild, { recursive: true, force: true });
    }

    // Create necessary directories
    [config.jsOutput, config.tailwindOutput, config.publicBuild].forEach(dir => {
        if (!existsSync(dir)) mkdirSync(dir, { recursive: true });
    });

    // Initial build
    await buildTypeScript(config.tsInput, isBuildMode ? config.jsBuild : config.jsOutput);
    await buildTailwind(config.tailwindInput, isBuildMode ? config.tailwindBuild : config.tailwindOutput);
    copyPublic(config.publicInput, config.publicBuild, config.publicIgnore);

    if (isBuildMode) {
        await buildGo();
        process.exit(0);
    }

    if (isWatchMode) {
        let goProcess: ChildProcess | null = runGoDev();
        let restartTimeout: NodeJS.Timeout | null = null;

        // Función para reiniciar el servidor Go
        const restartGoServer = async () => {
            if (restartTimeout) {
                clearTimeout(restartTimeout);
            }
            
            restartTimeout = setTimeout(async () => {
                if (goProcess) {
                    console.log('Reiniciando servidor Go...');
                    await stopGoProcess(goProcess);
                    // Espera un poquito por que el shutdown es asincrono
                    setTimeout(() => {
                        goProcess = runGoDev();
                    }, 500);
                }
                restartTimeout = null;
            }, 2000); // Debounce para que vscode no trolee, se puede ajuztar al poder del pc
        };

        // Watch for changes
        watch(config.tsInput, { recursive: true }, async (event, filename) => {
            if (filename && filename.endsWith('.ts')) {
                await buildTypeScript(config.tsInput, config.jsOutput);
            }
        });

        watch(config.tailwindInput, { recursive: true }, async (event, filename) => {
            if (filename && filename.endsWith('.css')) {
                await buildTailwind(config.tailwindInput, config.tailwindOutput);
            }
        });

        watch(config.publicInput, { recursive: true }, (event, filename) => {
            if (filename && !config.publicIgnore.some(ignore => filename.startsWith(ignore))) {
                copyPublic(config.publicInput, config.publicBuild, config.publicIgnore);
            }
        });

        // Watch para archivos Go
        const goFiles = findGoFiles(config.goInput);
            for (const goFile of goFiles) {
            watch(goFile, (event, filename) => {
                console.log(`Go file changed: ${goFile}`);
                restartGoServer();
            });
        }

        // También observa el archivo principal de Go
        watch(config.goMain, (event, filename) => {
            console.log('Main Go file changed');
            restartGoServer();
        });

        // Manejar cierre graceful del script
        process.on('SIGINT', async () => {
            console.log('Cerrando servidor Go...');
            if (goProcess) {
                await stopGoProcess(goProcess);
            }
            process.exit(0);
        });
    }
}

main().catch(console.error);
