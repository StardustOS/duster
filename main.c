#include <stdio.h>
#include <xenctrl.h>
#include <errno.h>
#include <stdio.h>
#include <fcntl.h>
#include <stdbool.h>

int main(void) {
    xc_interface* k = xc_interface_open(NULL, NULL, O_RDWR);
    xc_domain_pause(k, 26);
    xc_domain_setdebugging(k, 26, true);
    vcpu_guest_context_any_t m;
    int err = xc_vcpu_getcontext(k, 26, 0, &m);
    printf("%d\n", err);
    printf("RIP: %ld\n", m.x64.user_regs.rip);
    xc_domain_setdebugging(k, 26, false);
    xc_domain_unpause(k, 26);

    return 0;
}