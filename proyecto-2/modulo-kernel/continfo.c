#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/mm.h>
#include <linux/sched/signal.h>
#include <linux/slab.h>
#include <linux/string.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Nelson Cun");
MODULE_DESCRIPTION("Módulo que lista procesos de contenedores");
MODULE_VERSION("1.0");

#define PROC_NAME "continfo_so1_201222010"
#define MAX_CMDLINE_LENGTH 256

static char *get_process_cmdline(struct task_struct *task) {
    struct mm_struct *mm;
    char *cmdline;
    unsigned long arg_start, arg_end;
    int len, i;

    cmdline = kmalloc(MAX_CMDLINE_LENGTH, GFP_KERNEL);
    if (!cmdline) return NULL;

    mm = get_task_mm(task);
    if (!mm) { kfree(cmdline); return NULL; }

    down_read(&mm->mmap_lock);
    arg_start = mm->arg_start;
    arg_end = mm->arg_end;
    up_read(&mm->mmap_lock);

    len = arg_end - arg_start;
    if (len > MAX_CMDLINE_LENGTH - 1) len = MAX_CMDLINE_LENGTH - 1;

    if (access_process_vm(task, arg_start, cmdline, len, 0) != len) {
        mmput(mm);
        kfree(cmdline);
        return NULL;
    }

    cmdline[len] = '\0';
    for (i = 0; i < len; i++)
        if (cmdline[i] == '\0') cmdline[i] = ' ';

    mmput(mm);
    return cmdline;
}

static int continfo_show(struct seq_file *m, void *v) {
    struct task_struct *task;
    struct sysinfo si;

    si_meminfo(&si);
    seq_printf(m, "TotalRAM: %lu KB\n", si.totalram << (PAGE_SHIFT - 10));
    seq_printf(m, "FreeRAM: %lu KB\n\n", si.freeram << (PAGE_SHIFT - 10));
    seq_printf(m, "Procesos de Contenedores:\n");
    seq_printf(m, "PID\tNombre\tVSZ(KB)\tRSS(KB)\t%%Mem\t%%CPU\tID/Comando\n");

    for_each_process(task) {
        // Filtrar procesos que tengan "container" en comm o cmdline
        if (strstr(task->comm, "container") || strstr(get_process_cmdline(task), "container")) {
            unsigned long vsz = 0, rss = 0, mem_usage = 0, cpu_usage = 0;
            char *cmdline = get_process_cmdline(task);

            if (task->mm) {
                vsz = task->mm->total_vm << (PAGE_SHIFT - 10);
                rss = get_mm_rss(task->mm) << (PAGE_SHIFT - 10);
                mem_usage = (rss * 100) / (si.totalram << (PAGE_SHIFT - 10));
            }

            cpu_usage = 0; // simplificado

            seq_printf(m, "%d\t%s\t%lu\t%lu\t%lu\t%lu\t%s\n",
                       task->pid, task->comm, vsz, rss, mem_usage, cpu_usage,
                       cmdline ? cmdline : "N/A");

            if (cmdline) kfree(cmdline);
        }
    }

    return 0;
}

static int continfo_open(struct inode *inode, struct file *file) {
    return single_open(file, continfo_show, NULL);
}

static const struct proc_ops continfo_ops = {
    .proc_open = continfo_open,
    .proc_read = seq_read,
};

static int __init continfo_init(void) {
    if (!proc_create(PROC_NAME, 0, NULL, &continfo_ops)) {
        printk(KERN_ERR "No se pudo crear /proc/%s\n", PROC_NAME);
        return -ENOMEM;
    }
    printk(KERN_INFO "continfo modulo cargado\n");
    return 0;
}

static void __exit continfo_exit(void) {
    remove_proc_entry(PROC_NAME, NULL);
    printk(KERN_INFO "continfo modulo descargado\n");
}

module_init(continfo_init);
module_exit(continfo_exit);
