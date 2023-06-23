package provider

import (
	"telusag/terraform-provider-solace/sempv2"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

// MsgVpnQueue struct for MsgVpnQueue
type MsgVpnQueue struct {
	AccessType                              *string         `tfsdk:"access_type"`
	ConsumerAckPropagationEnabled           *bool           `tfsdk:"consumer_ack_propagation_enabled"`
	DeadMsgQueue                            *string         `tfsdk:"dead_msg_queue"`
	DeliveryCountEnabled                    *bool           `tfsdk:"delivery_count_enabled"`
	DeliveryDelay                           *int64          `tfsdk:"delivery_delay"`
	EgressEnabled                           *bool           `tfsdk:"egress_enabled"`
	EventBindCountThreshold                 *EventThreshold `tfsdk:"event_bind_count_threshold"`
	EventMsgSpoolUsageThreshold             *EventThreshold `tfsdk:"event_msg_spool_usage_threshold"`
	EventRejectLowPriorityMsgLimitThreshold *EventThreshold `tfsdk:"event_reject_low_priority_msg_limit_threshold"`
	IngressEnabled                          *bool           `tfsdk:"ingress_enabled"`
	MaxBindCount                            *int64          `tfsdk:"max_bind_count"`
	MaxDeliveredUnackedMsgsPerFlow          *int64          `tfsdk:"max_delivered_unacked_msgs_per_flow"`
	MaxMsgSize                              *int32          `tfsdk:"max_msg_size"`
	MaxMsgSpoolUsage                        *int64          `tfsdk:"max_msg_spool_usage"`
	MaxRedeliveryCount                      *int64          `tfsdk:"max_redelivery_count"`
	MaxTtl                                  *int64          `tfsdk:"max_ttl"`
	MsgVpnName                              *string         `tfsdk:"msg_vpn_name"`
	Owner                                   *string         `tfsdk:"owner"`
	Permission                              *string         `tfsdk:"permission"`
	QueueName                               *string         `tfsdk:"queue_name"`
	RedeliveryEnabled                       *bool           `tfsdk:"redelivery_enabled"`
	RejectLowPriorityMsgEnabled             *bool           `tfsdk:"reject_low_priority_msg_enabled"`
	RejectLowPriorityMsgLimit               *int64          `tfsdk:"reject_low_priority_msg_limit"`
	RejectMsgToSenderOnDiscardBehavior      *string         `tfsdk:"reject_msg_to_sender_on_discard_behavior"`
	RespectMsgPriorityEnabled               *bool           `tfsdk:"respect_msg_priority_enabled"`
	RespectTtlEnabled                       *bool           `tfsdk:"respect_ttl_enabled"`
}

func (tfData *MsgVpnQueue) ToTF(apiData *sempv2.MsgVpnQueue) {
	AssignIfDstNotNil(&tfData.AccessType, apiData.AccessType)
	AssignIfDstNotNil(&tfData.ConsumerAckPropagationEnabled, apiData.ConsumerAckPropagationEnabled)
	AssignIfDstNotNil(&tfData.DeadMsgQueue, apiData.DeadMsgQueue)
	AssignIfDstNotNil(&tfData.DeliveryCountEnabled, apiData.DeliveryCountEnabled)
	AssignIfDstNotNil(&tfData.DeliveryDelay, apiData.DeliveryDelay)
	AssignIfDstNotNil(&tfData.EgressEnabled, apiData.EgressEnabled)
	AssignIfDstNotNil(&tfData.EventBindCountThreshold, EventThresholdToTF(apiData.EventBindCountThreshold))
	AssignIfDstNotNil(&tfData.EventMsgSpoolUsageThreshold, EventThresholdToTF(apiData.EventMsgSpoolUsageThreshold))
	AssignIfDstNotNil(&tfData.EventRejectLowPriorityMsgLimitThreshold, EventThresholdToTF(apiData.EventRejectLowPriorityMsgLimitThreshold))
	AssignIfDstNotNil(&tfData.IngressEnabled, apiData.IngressEnabled)
	AssignIfDstNotNil(&tfData.MaxBindCount, apiData.MaxBindCount)
	AssignIfDstNotNil(&tfData.MaxDeliveredUnackedMsgsPerFlow, apiData.MaxDeliveredUnackedMsgsPerFlow)
	AssignIfDstNotNil(&tfData.MaxMsgSize, apiData.MaxMsgSize)
	AssignIfDstNotNil(&tfData.MaxMsgSpoolUsage, apiData.MaxMsgSpoolUsage)
	AssignIfDstNotNil(&tfData.MaxRedeliveryCount, apiData.MaxRedeliveryCount)
	AssignIfDstNotNil(&tfData.MaxTtl, apiData.MaxTtl)
	AssignIfDstNotNil(&tfData.MsgVpnName, apiData.MsgVpnName)
	AssignIfDstNotNil(&tfData.Owner, apiData.Owner)
	AssignIfDstNotNil(&tfData.Permission, apiData.Permission)
	AssignIfDstNotNil(&tfData.QueueName, apiData.QueueName)
	AssignIfDstNotNil(&tfData.RedeliveryEnabled, apiData.RedeliveryEnabled)
	AssignIfDstNotNil(&tfData.RejectLowPriorityMsgEnabled, apiData.RejectLowPriorityMsgEnabled)
	AssignIfDstNotNil(&tfData.RejectLowPriorityMsgLimit, apiData.RejectLowPriorityMsgLimit)
	AssignIfDstNotNil(&tfData.RejectMsgToSenderOnDiscardBehavior, apiData.RejectMsgToSenderOnDiscardBehavior)
	AssignIfDstNotNil(&tfData.RespectMsgPriorityEnabled, apiData.RespectMsgPriorityEnabled)
	AssignIfDstNotNil(&tfData.RespectTtlEnabled, apiData.RespectTtlEnabled)
}

func (tfData *MsgVpnQueue) ToApi() *sempv2.MsgVpnQueue {
	return &sempv2.MsgVpnQueue{
		AccessType:                              tfData.AccessType,
		ConsumerAckPropagationEnabled:           tfData.ConsumerAckPropagationEnabled,
		DeadMsgQueue:                            tfData.DeadMsgQueue,
		DeliveryCountEnabled:                    tfData.DeliveryCountEnabled,
		DeliveryDelay:                           tfData.DeliveryDelay,
		EgressEnabled:                           tfData.EgressEnabled,
		EventBindCountThreshold:                 tfData.EventBindCountThreshold.ToApi(),
		EventMsgSpoolUsageThreshold:             tfData.EventMsgSpoolUsageThreshold.ToApi(),
		EventRejectLowPriorityMsgLimitThreshold: tfData.EventRejectLowPriorityMsgLimitThreshold.ToApi(),
		IngressEnabled:                          tfData.IngressEnabled,
		MaxBindCount:                            tfData.MaxBindCount,
		MaxDeliveredUnackedMsgsPerFlow:          tfData.MaxDeliveredUnackedMsgsPerFlow,
		MaxMsgSize:                              tfData.MaxMsgSize,
		MaxMsgSpoolUsage:                        tfData.MaxMsgSpoolUsage,
		MaxRedeliveryCount:                      tfData.MaxRedeliveryCount,
		MaxTtl:                                  tfData.MaxTtl,
		MsgVpnName:                              tfData.MsgVpnName,
		Owner:                                   tfData.Owner,
		Permission:                              tfData.Permission,
		QueueName:                               tfData.QueueName,
		RedeliveryEnabled:                       tfData.RedeliveryEnabled,
		RejectLowPriorityMsgEnabled:             tfData.RejectLowPriorityMsgEnabled,
		RejectLowPriorityMsgLimit:               tfData.RejectLowPriorityMsgLimit,
		RejectMsgToSenderOnDiscardBehavior:      tfData.RejectMsgToSenderOnDiscardBehavior,
		RespectMsgPriorityEnabled:               tfData.RespectMsgPriorityEnabled,
		RespectTtlEnabled:                       tfData.RespectTtlEnabled,
	}
}

// Terraform Resource schema for MsgVpnQueue
func MsgVpnQueueResourceSchema(requiredAttributes ...string) schema.Schema {
	schema := schema.Schema{
		Description: "MsgVpnQueue",
		Attributes: map[string]schema.Attribute{
			"access_type": schema.StringAttribute{
				Description: "The access type for delivering messages to consumer flows bound to the Queue. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"exclusive\"`. The allowed values and their meaning are:  <pre> \"exclusive\" - Exclusive delivery of messages to the first bound consumer flow. \"non-exclusive\" - Non-exclusive delivery of messages to all bound consumer flows in a round-robin fashion. </pre> ",
				Required:    contains(requiredAttributes, "access_type"),
				Optional:    !contains(requiredAttributes, "access_type"),

				Validators: []validator.String{
					stringvalidator.OneOf("exclusive", "non-exclusive"),
				},
				PlanModifiers: StringPlanModifiersFor("access_type", requiredAttributes),
			},
			"consumer_ack_propagation_enabled": schema.BoolAttribute{
				Description: "Enable or disable the propagation of consumer acknowledgements (ACKs) received on the active replication Message VPN to the standby replication Message VPN. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`.",
				Required:    contains(requiredAttributes, "consumer_ack_propagation_enabled"),
				Optional:    !contains(requiredAttributes, "consumer_ack_propagation_enabled"),
			},
			"dead_msg_queue": schema.StringAttribute{
				Description: "The name of the Dead Message Queue (DMQ) used by the Queue. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"#DEAD_MSG_QUEUE\"`. Available since 2.2.",
				Required:    contains(requiredAttributes, "dead_msg_queue"),
				Optional:    !contains(requiredAttributes, "dead_msg_queue"),

				PlanModifiers: StringPlanModifiersFor("dead_msg_queue", requiredAttributes),
			},
			"delivery_count_enabled": schema.BoolAttribute{
				Description: "Enable or disable the ability for client applications to query the message delivery count of messages received from the Queue. This is a controlled availability feature. Please contact support to find out if this feature is supported for your use case. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.19.",
				Required:    contains(requiredAttributes, "delivery_count_enabled"),
				Optional:    !contains(requiredAttributes, "delivery_count_enabled"),
			},
			"delivery_delay": schema.Int64Attribute{
				Description: "The delay, in seconds, to apply to messages arriving on the Queue before the messages are eligible for delivery. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`. Available since 2.22.",
				Required:    contains(requiredAttributes, "delivery_delay"),
				Optional:    !contains(requiredAttributes, "delivery_delay"),
			},
			"egress_enabled": schema.BoolAttribute{
				Description: "Enable or disable the transmission of messages from the Queue. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "egress_enabled"),
				Optional:    !contains(requiredAttributes, "egress_enabled"),
			},
			"event_bind_count_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_bind_count_threshold"),
				Optional:    !contains(requiredAttributes, "event_bind_count_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"event_msg_spool_usage_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_msg_spool_usage_threshold"),
				Optional:    !contains(requiredAttributes, "event_msg_spool_usage_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"event_reject_low_priority_msg_limit_threshold": schema.SingleNestedAttribute{
				Description: "",
				Required:    contains(requiredAttributes, "event_reject_low_priority_msg_limit_threshold"),
				Optional:    !contains(requiredAttributes, "event_reject_low_priority_msg_limit_threshold"),

				Attributes: EventThresholdResourceAttributes,
			},
			"ingress_enabled": schema.BoolAttribute{
				Description: "Enable or disable the reception of messages to the Queue. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "ingress_enabled"),
				Optional:    !contains(requiredAttributes, "ingress_enabled"),
			},
			"max_bind_count": schema.Int64Attribute{
				Description: "The maximum number of consumer flows that can bind to the Queue. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `1000`.",
				Required:    contains(requiredAttributes, "max_bind_count"),
				Optional:    !contains(requiredAttributes, "max_bind_count"),
			},
			"max_delivered_unacked_msgs_per_flow": schema.Int64Attribute{
				Description: "The maximum number of messages delivered but not acknowledged per flow for the Queue. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `10000`.",
				Required:    contains(requiredAttributes, "max_delivered_unacked_msgs_per_flow"),
				Optional:    !contains(requiredAttributes, "max_delivered_unacked_msgs_per_flow"),
			},
			"max_msg_size": schema.Int64Attribute{
				Description: "The maximum message size allowed in the Queue, in bytes (B). Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `10000000`.",
				Required:    contains(requiredAttributes, "max_msg_size"),
				Optional:    !contains(requiredAttributes, "max_msg_size"),
			},
			"max_msg_spool_usage": schema.Int64Attribute{
				Description: "The maximum message spool usage allowed by the Queue, in megabytes (MB). A value of 0 only allows spooling of the last message received and disables quota checking. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `5000`.",
				Required:    contains(requiredAttributes, "max_msg_spool_usage"),
				Optional:    !contains(requiredAttributes, "max_msg_spool_usage"),
			},
			"max_redelivery_count": schema.Int64Attribute{
				Description: "The maximum number of times the Queue will attempt redelivery of a message prior to it being discarded or moved to the DMQ. A value of 0 means to retry forever. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.",
				Required:    contains(requiredAttributes, "max_redelivery_count"),
				Optional:    !contains(requiredAttributes, "max_redelivery_count"),
			},
			"max_ttl": schema.Int64Attribute{
				Description: "The maximum time in seconds a message can stay in the Queue when `respectTtlEnabled` is `\"true\"`. A message expires when the lesser of the sender assigned time-to-live (TTL) in the message and the `maxTtl` configured for the Queue, is exceeded. A value of 0 disables expiry. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.",
				Required:    contains(requiredAttributes, "max_ttl"),
				Optional:    !contains(requiredAttributes, "max_ttl"),
			},
			"msg_vpn_name": schema.StringAttribute{
				Description: "The name of the Message VPN.",
				Required:    contains(requiredAttributes, "msg_vpn_name"),
				Optional:    !contains(requiredAttributes, "msg_vpn_name"),

				PlanModifiers: StringPlanModifiersFor("msg_vpn_name", requiredAttributes),
			},
			"owner": schema.StringAttribute{
				Description: "The Client Username that owns the Queue and has permission equivalent to `\"delete\"`. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"\"`.",
				Required:    contains(requiredAttributes, "owner"),
				Optional:    !contains(requiredAttributes, "owner"),

				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[[:print:]]{1,189}$"), "Does not match pattern '^[[:print:]]{1,189}$'"),
				},
				PlanModifiers: StringPlanModifiersFor("owner", requiredAttributes),
			},
			"permission": schema.StringAttribute{
				Description: "The permission level for all consumers of the Queue, excluding the owner. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"no-access\"`. The allowed values and their meaning are:  <pre> \"no-access\" - Disallows all access. \"read-only\" - Read-only access to the messages. \"consume\" - Consume (read and remove) messages. \"modify-topic\" - Consume messages or modify the topic/selector. \"delete\" - Consume messages, modify the topic/selector or delete the Client created endpoint altogether. </pre> ",
				Required:    contains(requiredAttributes, "permission"),
				Optional:    !contains(requiredAttributes, "permission"),

				Validators: []validator.String{
					stringvalidator.OneOf("no-access", "read-only", "consume", "modify-topic", "delete"),
				},
				PlanModifiers: StringPlanModifiersFor("permission", requiredAttributes),
			},
			"queue_name": schema.StringAttribute{
				Description: "The name of the Queue.",
				Required:    contains(requiredAttributes, "queue_name"),
				Optional:    !contains(requiredAttributes, "queue_name"),

				PlanModifiers: StringPlanModifiersFor("queue_name", requiredAttributes),
			},
			"redelivery_enabled": schema.BoolAttribute{
				Description: "Enable or disable message redelivery. When enabled, the number of redelivery attempts is controlled by maxRedeliveryCount. When disabled, the message will never be delivered from the queue more than once. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `true`. Available since 2.18.",
				Required:    contains(requiredAttributes, "redelivery_enabled"),
				Optional:    !contains(requiredAttributes, "redelivery_enabled"),
			},
			"reject_low_priority_msg_enabled": schema.BoolAttribute{
				Description: "Enable or disable the checking of low priority messages against the `rejectLowPriorityMsgLimit`. This may only be enabled if `rejectMsgToSenderOnDiscardBehavior` does not have a value of `\"never\"`. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "reject_low_priority_msg_enabled"),
				Optional:    !contains(requiredAttributes, "reject_low_priority_msg_enabled"),
			},
			"reject_low_priority_msg_limit": schema.Int64Attribute{
				Description: "The number of messages of any priority in the Queue above which low priority messages are not admitted but higher priority messages are allowed. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `0`.",
				Required:    contains(requiredAttributes, "reject_low_priority_msg_limit"),
				Optional:    !contains(requiredAttributes, "reject_low_priority_msg_limit"),
			},
			"reject_msg_to_sender_on_discard_behavior": schema.StringAttribute{
				Description: "Determines when to return negative acknowledgements (NACKs) to sending clients on message discards. Note that NACKs cause the message to not be delivered to any destination and Transacted Session commits to fail. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `\"when-queue-enabled\"`. The allowed values and their meaning are:  <pre> \"always\" - Always return a negative acknowledgment (NACK) to the sending client on message discard. \"when-queue-enabled\" - Only return a negative acknowledgment (NACK) to the sending client on message discard when the Queue is enabled. \"never\" - Never return a negative acknowledgment (NACK) to the sending client on message discard. </pre>  Available since 2.1.",
				Required:    contains(requiredAttributes, "reject_msg_to_sender_on_discard_behavior"),
				Optional:    !contains(requiredAttributes, "reject_msg_to_sender_on_discard_behavior"),

				Validators: []validator.String{
					stringvalidator.OneOf("always", "when-queue-enabled", "never"),
				},
				PlanModifiers: StringPlanModifiersFor("reject_msg_to_sender_on_discard_behavior", requiredAttributes),
			},
			"respect_msg_priority_enabled": schema.BoolAttribute{
				Description: "Enable or disable the respecting of message priority. When enabled, messages contained in the Queue are delivered in priority order, from 9 (highest) to 0 (lowest). MQTT queues do not support enabling message priority. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`. Available since 2.8.",
				Required:    contains(requiredAttributes, "respect_msg_priority_enabled"),
				Optional:    !contains(requiredAttributes, "respect_msg_priority_enabled"),
			},
			"respect_ttl_enabled": schema.BoolAttribute{
				Description: "Enable or disable the respecting of the time-to-live (TTL) for messages in the Queue. When enabled, expired messages are discarded or moved to the DMQ. Changes to this attribute are synchronized to HA mates and replication sites via config-sync. The default value is `false`.",
				Required:    contains(requiredAttributes, "respect_ttl_enabled"),
				Optional:    !contains(requiredAttributes, "respect_ttl_enabled"),
			},
		},
	}

	return schema
}
