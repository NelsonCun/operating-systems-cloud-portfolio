#include <linux/module.h>
#include <linux/export-internal.h>
#include <linux/compiler.h>

MODULE_INFO(name, KBUILD_MODNAME);

__visible struct module __this_module
__section(".gnu.linkonce.this_module") = {
	.name = KBUILD_MODNAME,
	.init = init_module,
#ifdef CONFIG_MODULE_UNLOAD
	.exit = cleanup_module,
#endif
	.arch = MODULE_ARCH_INIT,
};



static const struct modversion_info ____versions[]
__used __section("__versions") = {
	{ 0x5218fe90, "single_open" },
	{ 0xbd03ed67, "random_kmalloc_seed" },
	{ 0x4ac4312d, "kmalloc_caches" },
	{ 0x8d1d7639, "__kmalloc_cache_noprof" },
	{ 0xe3c29035, "get_task_mm" },
	{ 0xa59da3c0, "down_read" },
	{ 0xa59da3c0, "up_read" },
	{ 0x013f9a5f, "access_process_vm" },
	{ 0x0142c002, "mmput" },
	{ 0xcb8b6ec6, "kfree" },
	{ 0xc7ffe1aa, "si_meminfo" },
	{ 0x12cfb334, "seq_printf" },
	{ 0xa2152099, "init_task" },
	{ 0x17545440, "strstr" },
	{ 0xd272d446, "__stack_chk_fail" },
	{ 0xfefac423, "remove_proc_entry" },
	{ 0xd22cd56f, "seq_read" },
	{ 0xd272d446, "__fentry__" },
	{ 0xf8d7ac5e, "proc_create" },
	{ 0xe8213e80, "_printk" },
	{ 0xd272d446, "__x86_return_thunk" },
	{ 0x70eca2ca, "module_layout" },
};

static const u32 ____version_ext_crcs[]
__used __section("__version_ext_crcs") = {
	0x5218fe90,
	0xbd03ed67,
	0x4ac4312d,
	0x8d1d7639,
	0xe3c29035,
	0xa59da3c0,
	0xa59da3c0,
	0x013f9a5f,
	0x0142c002,
	0xcb8b6ec6,
	0xc7ffe1aa,
	0x12cfb334,
	0xa2152099,
	0x17545440,
	0xd272d446,
	0xfefac423,
	0xd22cd56f,
	0xd272d446,
	0xf8d7ac5e,
	0xe8213e80,
	0xd272d446,
	0x70eca2ca,
};
static const char ____version_ext_names[]
__used __section("__version_ext_names") =
	"single_open\0"
	"random_kmalloc_seed\0"
	"kmalloc_caches\0"
	"__kmalloc_cache_noprof\0"
	"get_task_mm\0"
	"down_read\0"
	"up_read\0"
	"access_process_vm\0"
	"mmput\0"
	"kfree\0"
	"si_meminfo\0"
	"seq_printf\0"
	"init_task\0"
	"strstr\0"
	"__stack_chk_fail\0"
	"remove_proc_entry\0"
	"seq_read\0"
	"__fentry__\0"
	"proc_create\0"
	"_printk\0"
	"__x86_return_thunk\0"
	"module_layout\0"
;

MODULE_INFO(depends, "");


MODULE_INFO(srcversion, "7F34B4A3B0E0E12769CBE15");
