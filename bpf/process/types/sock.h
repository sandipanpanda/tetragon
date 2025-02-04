// SPDX-License-Identifier: GPL-2.0
/* Copyright Authors of Cilium */

#ifndef __SOCK_H__
#define __SOCK_H__

#include "tuple.h"

struct sk_type {
	struct tuple_type tuple;
	__u32 mark;
	__u32 priority;
	__u16 family;
	__u16 type;
	__u32 pad;
};

/* set_event_from_sock(sock)
 *
 * Populate the event args with the sock info.
 */
static inline __attribute__((unused)) void
set_event_from_sock(struct sk_type *event, struct sock *sk)
{
	struct sock_common *common = (struct sock_common *)sk;

	event->family = 0;

	probe_read(&event->family, sizeof(event->family),
		   _(&common->skc_family));
	probe_read(&event->type, sizeof(event->type), _(&sk->sk_type));
	probe_read(&event->tuple.protocol, sizeof(event->tuple.protocol),
		   _(&sk->sk_protocol));
	if (bpf_core_field_size(sk->sk_protocol) == 4) {
		// In the BTF, the protocol field in kernels <v5.6 is 8 bits of a u32.
		// As such, the easiest way to get the correct (8 bit) value is to read
		// it as a u16 and shift it by 1 byte.
		event->tuple.protocol = event->tuple.protocol >> 8;
	}
	probe_read(&event->mark, sizeof(event->mark), _(&sk->sk_mark));
	probe_read(&event->priority, sizeof(event->priority),
		   _(&sk->sk_priority));

	probe_read(&event->tuple.saddr, sizeof(event->tuple.saddr), _(&common->skc_rcv_saddr));
	probe_read(&event->tuple.daddr, sizeof(event->tuple.daddr), _(&common->skc_daddr));
	probe_read(&event->tuple.sport, sizeof(event->tuple.sport), _(&common->skc_num));
	probe_read(&event->tuple.dport, sizeof(event->tuple.dport), _(&common->skc_dport));
	event->tuple.dport = bpf_ntohs(event->tuple.dport);
}
#endif // __SOCK_H__
