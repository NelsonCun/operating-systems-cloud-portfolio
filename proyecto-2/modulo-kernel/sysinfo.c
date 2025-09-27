#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/mm.h>
#include <linux/sched/signal.h>
#include <linux/slab.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Nelson Cun");
MODULE_DESCRIPTION("Módulo que lista todos los procesos del sistema");
MODULE_VERSION("1.0");

#define PROC_NAME "sysinfo_so1_201222010"
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

static int sysinfo_show(struct seq_file *m, void *v) {
    struct task_struct *task;
    struct sysinfo si;

    si_meminfo(&si);
    seq_printf(m, "TotalRAM: %lu KB\n", si.totalram << (PAGE_SHIFT - 10));
    seq_printf(m, "FreeRAM: %lu KB\n\n", si.freeram << (PAGE_SHIFT - 10));
    seq_printf(m, "Procesos del Sistema:\n");
    seq_printf(m, "PID\tNombre\tVSZ(KB)\tRSS(KB)\t%%Mem\t%%CPU\tEstado\tCmdline\n");

    for_each_process(task) {
        unsigned long vsz = 0, rss = 0, mem_usage = 0, cpu_usage = 0;
        char *cmdline = get_process_cmdline(task);
        char state_str[16];

        // Estado del proceso
        if (task_is_running(task))
            snprintf(state_str, sizeof(state_str), "RUNNING");
        else if (task->flags & PF_EXITING)
            snprintf(state_str, sizeof(state_str), "EXITING");
        else
            snprintf(state_str, sizeof(state_str), "SLEEPING");


        // Memoria
        if (task->mm) {
            vsz = task->mm->total_vm << (PAGE_SHIFT - 10);
            rss = get_mm_rss(task->mm) << (PAGE_SHIFT - 10);
            mem_usage = (rss * 100) / (si.totalram << (PAGE_SHIFT - 10));
        }

        // CPU (aproximado)
        cpu_usage = 0; // Opcional simplificado, ya que calcular %CPU real requiere más

        seq_printf(m, "%d\t%s\t%lu\t%lu\t%lu\t%lu\t%s\t%s\n",
                   task->pid, task->comm, vsz, rss, mem_usage, cpu_usage,
                   state_str, cmdline ? cmdline : "N/A");

        if (cmdline) kfree(cmdline);
    }

    return 0;
}

static int sysinfo_open(struct inode *inode, struct file *file) {
    return single_open(file, sysinfo_show, NULL);
}

static const struct proc_ops sysinfo_ops = {
    .proc_open = sysinfo_open,
    .proc_read = seq_read,
};

static int __init sysinfo_init(void) {
    if (!proc_create(PROC_NAME, 0, NULL, &sysinfo_ops)) {
        printk(KERN_ERR "No se pudo crear /proc/%s\n", PROC_NAME);
        return -ENOMEM;
    }
    printk(KERN_INFO "sysinfo modulo cargado\n");
    return 0;
}

static void __exit sysinfo_exit(void) {
    remove_proc_entry(PROC_NAME, NULL);
    printk(KERN_INFO "sysinfo modulo descargado\n");
}

module_init(sysinfo_init);
module_exit(sysinfo_exit);
